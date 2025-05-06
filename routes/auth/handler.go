package auth

import (
	"log" // Changed from fmt
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yantology/simple-pos/config"
	"github.com/yantology/simple-pos/pkg/dto"
	"github.com/yantology/simple-pos/pkg/resendutils"
)

type authHandler struct {
	authService    AuthService
	authRepository *AuthRepository
	emailSender    resendutils.ResendUtilsInterface
	emailTemplate  EmailTemplateInterface
	tokenRequest   *config.TokenConfig
}

func NewAuthHandler(
	authService AuthService,
	authRepository *AuthRepository,
	emailSender resendutils.ResendUtilsInterface,
	emailTemplate EmailTemplateInterface,
	tokenRequest *config.TokenConfig,
) *authHandler {
	return &authHandler{
		authService:    authService,
		authRepository: authRepository,
		emailSender:    emailSender,
		emailTemplate:  emailTemplate,
		tokenRequest:   tokenRequest,
	}
}

// @Summary Request activation token
// @Description Request a token for registration or password reset
// @Tags auth
// @Accept json
// @Produce json
// @Param type path string true "Token type (registration or forget-password)"
// @Param request body TokenRequest true "Token request parameters"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 404 {object} dto.MessageResponse
// @Failure 409 {object} dto.MessageResponse
// @Router /auth/token/{type} [post]
func (h *authHandler) RequestToken(c *gin.Context) {
	log.Printf("[AuthHandler] RequestToken: Started for type: %s, RequestID: %s\n", c.Param("type"), c.GetString("RequestID"))
	tokenType := c.Param("type")

	// Validate token type
	log.Printf("[AuthHandler] RequestToken: Validating token type: %s, RequestID: %s\n", tokenType, c.GetString("RequestID"))
	if tokenType != "registration" && tokenType != "forget-password" {
		log.Printf("[AuthHandler] RequestToken: Invalid token type: %s, RequestID: %s\n", tokenType, c.GetString("RequestID"))
		c.JSON(http.StatusBadRequest, dto.MessageResponse{
			Message: "Tipe token tidak valid",
		})
		return
	}

	var req TokenRequest
	if cuserr := c.ShouldBindJSON(&req); cuserr != nil {
		log.Printf("[AuthHandler] RequestToken: Invalid request format: %v, RequestID: %s\n", cuserr, c.GetString("RequestID"))
		c.JSON(http.StatusBadRequest, dto.MessageResponse{
			Message: "Format request tidak valid",
		})
		return
	}
	log.Printf("[AuthHandler] RequestToken: Received email: %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))

	// Validate email format
	log.Printf("[AuthHandler] RequestToken: Validating email: %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	if cuserr := h.authService.ValidateEmail(req.Email); cuserr != nil {
		log.Printf("[AuthHandler] RequestToken: Email validation failed for %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	// Check if email exists based on token type
	if tokenType == "registration" {
		log.Printf("[AuthHandler] RequestToken: Checking if email not already registered: %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
		if cuserr := h.authRepository.CheckIsNotExistingEmail(req.Email); cuserr != nil {
			log.Printf("[AuthHandler] RequestToken: Email %s already exists for registration, RequestID: %s\n", req.Email, c.GetString("RequestID"))
			c.JSON(http.StatusConflict, dto.MessageResponse{
				Message: "Email sudah terdaftar.", // Consistent message
			})
			return
		}
	} else if tokenType == "forget-password" {
		log.Printf("[AuthHandler] RequestToken: Checking if email exists for password reset: %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
		if cuserr := h.authRepository.CheckIsExistingEmail(req.Email); cuserr != nil {
			log.Printf("[AuthHandler] RequestToken: Email %s not found for password reset, RequestID: %s\n", req.Email, c.GetString("RequestID"))
			c.JSON(http.StatusNotFound, dto.MessageResponse{
				Message: cuserr.Message(),
			})
			return
		}
	}

	// Generate activation token
	log.Printf("[AuthHandler] RequestToken: Generating activation token for %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	token, cuserr := h.authService.GenerateActivationToken()
	if cuserr != nil {
		log.Printf("[AuthHandler] RequestToken: Failed to generate token for %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	// Hash the token before storing
	log.Printf("[AuthHandler] RequestToken: Hashing the token for %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	hashedToken, cuserr := h.authService.HashString(token)
	if cuserr != nil {
		log.Printf("[AuthHandler] RequestToken: Failed to hash token for %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	// Save token to database
	tokenReq := &ActivationTokenRequest{
		Email:          req.Email,
		ActivationCode: hashedToken,
		TokenType:      tokenType,
		ExpiryMinutes:  int(h.tokenRequest.AccessTokenExpiry.Minutes()), // Use AccessTokenExpiry converted to int minutes
	}

	log.Printf("[AuthHandler] RequestToken: Saving token to database for email: %s, type: %s, RequestID: %s\n", req.Email, tokenType, c.GetString("RequestID"))
	if cuserr := h.authRepository.SaveActivationToken(tokenReq); cuserr != nil {
		log.Printf("[AuthHandler] RequestToken: Failed to save token for %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: "Gagal menyimpan token",
		})
		return
	}

	// Generate email content based on token type
	var emailHTML, emailSubject string
	if tokenType == "registration" {
		log.Printf("[AuthHandler] RequestToken: Generating registration email for: %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
		emailHTML = h.emailTemplate.GenerateRegistrationEmail(req.Email, token)
		emailSubject = "Kode Aktivasi Pendaftaran Akun Anda"
	}
	if tokenType == "forget-password" {
		log.Printf("[AuthHandler] RequestToken: Generating password reset email for: %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
		emailHTML = h.emailTemplate.GeneratePasswordResetEmail(req.Email, token)
		emailSubject = "Kode Reset Password Akun Anda"
	}

	// Send email
	log.Printf("[AuthHandler] RequestToken: Sending email to: %s, subject: %s, RequestID: %s\n", req.Email, emailSubject, c.GetString("RequestID"))
	if cuserr := h.emailSender.Send(emailHTML, emailSubject, []string{req.Email}); cuserr != nil {
		log.Printf("[AuthHandler] RequestToken: Failed to send email to %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		// Do not return error to client if email sending fails, but log it.
		// The token is already saved. User can retry or use "resend token" if implemented.
		// However, for critical actions like registration, you might want to inform the user.
		// For now, let's proceed and let the user know the token was "sent" (or attempted).
	}

	log.Printf("[AuthHandler] RequestToken: Successfully processed token request for email: %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Kode aktivasi telah dikirim ke email Anda. Silakan periksa folder inbox atau spam.",
	})
}

// @Summary Register new user
// @Description Register a new user with activation code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} dto.MessageResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 401 {object} dto.MessageResponse
// @Router /auth/register [post]
func (h *authHandler) Register(c *gin.Context) {
	log.Printf("[AuthHandler] Register: Started, RequestID: %s\n", c.GetString("RequestID"))
	var req RegisterRequest
	if cuserr := c.ShouldBindJSON(&req); cuserr != nil {
		log.Printf("[AuthHandler] Register: Invalid request format: %v, RequestID: %s\n", cuserr, c.GetString("RequestID"))
		c.JSON(http.StatusBadRequest, dto.MessageResponse{
			Message: "Format request tidak valid",
		})
		return
	}
	log.Printf("[AuthHandler] Register: Received registration request for email: %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))

	// Validate registration input
	regReq := RegistrationRequest{
		Email:                req.Email,
		Username:             req.Fullname, // Assuming Fullname is used as Username here
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	}

	log.Printf("[AuthHandler] Register: Validating registration input for %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	if cuserr := h.authService.ValidateRegistrationInput(regReq); cuserr != nil {
		log.Printf("[AuthHandler] Register: Registration input validation failed for %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	// Validate activation token
	tokenReq := &GetActivationTokenRequest{
		Email:     req.Email,
		TokenType: "registration",
	}

	log.Printf("[AuthHandler] Register: Validating activation token for %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	token, cuserr := h.authRepository.GetActivationToken(tokenReq)
	if cuserr != nil {
		log.Printf("[AuthHandler] Register: Failed to get activation token for %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	// Verify token
	log.Printf("[AuthHandler] Register: Verifying activation code for %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	if cuserr := h.authService.VerifyHash(token, req.ActivationCode); cuserr != nil {
		log.Printf("[AuthHandler] Register: Activation code verification failed for %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(), // "Kode aktivasi tidak valid atau sudah kedaluwarsa."
		})
		return
	}

	// Hash password
	log.Printf("[AuthHandler] Register: Hashing password for %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	hashedPassword, cuserr := h.authService.HashString(req.Password)
	if cuserr != nil {
		log.Printf("[AuthHandler] Register: Failed to hash password for %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	// Create user
	createUserReq := &CreateUserRequest{
		Email:        req.Email,
		Fullname:     req.Fullname,
		PasswordHash: hashedPassword,
	}

	log.Printf("[AuthHandler] Register: Creating user %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	if cuserr := h.authRepository.CreateUser(createUserReq); cuserr != nil {
		log.Printf("[AuthHandler] Register: Failed to create user %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(), // "Gagal membuat pengguna."
		})
		return
	}

	log.Printf("[AuthHandler] Register: User %s registered successfully, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	c.JSON(http.StatusCreated, dto.MessageResponse{
		Message: "Pendaftaran berhasil, silakan login."})
}

// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 401 {object} dto.MessageResponse
// @Router /auth/login [post]
func (h *authHandler) Login(c *gin.Context) {
	log.Printf("[AuthHandler] Login: Started, RequestID: %s\n", c.GetString("RequestID"))
	var req LoginRequest
	if cuserr := c.ShouldBindJSON(&req); cuserr != nil {
		log.Printf("[AuthHandler] Login: Invalid request format: %v, RequestID: %s\n", cuserr, c.GetString("RequestID"))
		c.JSON(http.StatusBadRequest, dto.MessageResponse{
			Message: "Format request tidak valid.",
		})
		return
	}
	log.Printf("[AuthHandler] Login: Received login request for email: %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))

	// Get user by email
	log.Printf("[AuthHandler] Login: Getting user by email: %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	user, cuserr := h.authRepository.GetUserByEmail(req.Email)
	if cuserr != nil {
		log.Printf("[AuthHandler] Login: Failed to get user %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{
			Message: "Email atau password salah.", // Generic message for security
		})
		return
	}

	// Verify password
	log.Printf("[AuthHandler] Login: Verifying password for user %s (ID: %d), RequestID: %s\n", user.Email, user.ID, c.GetString("RequestID"))
	if cuserr := h.authService.VerifyHash(user.PasswordHash, req.Password); cuserr != nil {
		log.Printf("[AuthHandler] Login: Password verification failed for user %s: %s, RequestID: %s\n", user.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(http.StatusUnauthorized, dto.MessageResponse{
			Message: "Email atau password salah.", // Generic message
		})
		return
	}

	// Generate token pair
	tokenPairReq := TokenPairRequest{
		UserID: user.ID,
		Email:  user.Email,
	}

	log.Printf("[AuthHandler] Login: Generating token pair for user %s (ID: %d), RequestID: %s\n", user.Email, user.ID, c.GetString("RequestID"))
	cuserr = h.authService.GenerateTokenPairCookies(c.Writer, tokenPairReq)
	if cuserr != nil {
		log.Printf("[AuthHandler] Login: Failed to generate token pair for user %s: %s, RequestID: %s\n", user.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	log.Printf("[AuthHandler] Login: User %s (ID: %d) logged in successfully, RequestID: %s\n", user.Email, user.ID, c.GetString("RequestID"))
	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Login berhasil.",
	})
}

// @Summary Reset password
// @Description Reset user password using activation code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ForgetPasswordRequest true "Password reset details"
// @Success 200 {object} dto.MessageResponse "Success response with message"
// @Failure 400 {object} dto.MessageResponse "Bad request response"
// @Failure 401 {object} dto.MessageResponse "Unauthorized response"
// @Router /auth/forget-password [post]
func (h *authHandler) ForgetPassword(c *gin.Context) {
	log.Printf("[AuthHandler] ForgetPassword: Started, RequestID: %s\n", c.GetString("RequestID"))
	var req ForgetPasswordRequest
	if cuserr := c.ShouldBindJSON(&req); cuserr != nil {
		log.Printf("[AuthHandler] ForgetPassword: Invalid request format: %v, RequestID: %s\n", cuserr, c.GetString("RequestID"))
		c.JSON(http.StatusBadRequest, dto.MessageResponse{
			Message: "Format request tidak valid.",
		})
		return
	}
	log.Printf("[AuthHandler] ForgetPassword: Received request for email: %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))

	// Validate password match
	log.Printf("[AuthHandler] ForgetPassword: Validating password input for %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	if cuserr := h.authService.ValidatePasswordInput(req.NewPassword, req.NewPasswordConfirmation); cuserr != nil {
		log.Printf("[AuthHandler] ForgetPassword: Password input validation failed for %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	// Validate activation token
	tokenReq := &GetActivationTokenRequest{
		Email:     req.Email,
		TokenType: "forget-password",
	}

	log.Printf("[AuthHandler] ForgetPassword: Validating activation token for %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	token, cuserr := h.authRepository.GetActivationToken(tokenReq)
	if cuserr != nil {
		log.Printf("[AuthHandler] ForgetPassword: Failed to get activation token for %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	// Verify token
	log.Printf("[AuthHandler] ForgetPassword: Verifying activation code for %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	if cuserr := h.authService.VerifyHash(token, req.ActivationCode); cuserr != nil {
		log.Printf("[AuthHandler] ForgetPassword: Activation code verification failed for %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(), // "Kode aktivasi tidak valid atau sudah kedaluwarsa."
		})
		return
	}

	// Hash new password
	log.Printf("[AuthHandler] ForgetPassword: Hashing new password for %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	hashedPassword, cuserr := h.authService.HashString(req.NewPassword)
	if cuserr != nil {
		log.Printf("[AuthHandler] ForgetPassword: Failed to hash new password for %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	// Update user password
	updateReq := &UpdatePasswordRequest{
		Email:           req.Email,
		NewPasswordHash: hashedPassword,
	}

	log.Printf("[AuthHandler] ForgetPassword: Updating password for user %s, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	if cuserr := h.authRepository.UpdateUserPassword(updateReq); cuserr != nil {
		log.Printf("[AuthHandler] ForgetPassword: Failed to update password for user %s: %s, RequestID: %s\n", req.Email, cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(), // "Gagal mengubah password."
		})
		return
	}

	log.Printf("[AuthHandler] ForgetPassword: Password for user %s reset successfully, RequestID: %s\n", req.Email, c.GetString("RequestID"))
	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Password berhasil diubah. Silakan login dengan password baru Anda.",
	})
}

// @Summary Refresh token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 401 {object} dto.MessageResponse
// @Router /auth/refresh-token [get]
func (h *authHandler) RefreshToken(c *gin.Context) {
	log.Printf("[AuthHandler] RefreshToken: Started, RequestID: %s\n", c.GetString("RequestID"))
	// Get refresh token from cookies
	refreshToken, err := c.Cookie(h.tokenRequest.RefreshTokenName)
	if err != nil {
		log.Printf("[AuthHandler] RefreshToken: Refresh token not found in cookies, RequestID: %s\n", c.GetString("RequestID"))
		c.JSON(http.StatusBadRequest, dto.MessageResponse{
			Message: "Refresh token tidak ditemukan dalam cookies.",
		})
		return
	}
	log.Printf("[AuthHandler] RefreshToken: Refresh token found in cookies, RequestID: %s\n", c.GetString("RequestID"))

	// Validate refresh token
	log.Printf("[AuthHandler] RefreshToken: Validating refresh token, RequestID: %s\n", c.GetString("RequestID"))
	claims, cuserr := h.authService.ValidateRefreshTokenClaims(refreshToken)
	if cuserr != nil {
		log.Printf("[AuthHandler] RefreshToken: Refresh token validation failed: %s, RequestID: %s\n", cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}
	log.Printf("[AuthHandler] RefreshToken: Refresh token validated for UserID: %s, Email: %s, RequestID: %s\n", claims.UserID, claims.Email, c.GetString("RequestID"))

	// Convert claims.UserID from string to int
	userID, err := strconv.Atoi(claims.UserID)
	if err != nil {
		log.Printf("[AuthHandler] RefreshToken: Invalid user ID format in claims '%s': %v, RequestID: %s\n", claims.UserID, err, c.GetString("RequestID"))
		c.JSON(http.StatusBadRequest, dto.MessageResponse{
			Message: "Invalid user ID format.",
		})
		return
	}

	tokenPair := TokenPairRequest{
		UserID: userID,
		Email:  claims.Email,
	}

	// Generate new access token
	log.Printf("[AuthHandler] RefreshToken: Generating new token pair for UserID: %d, Email: %s, RequestID: %s\n", userID, claims.Email, c.GetString("RequestID"))
	cuserr = h.authService.GenerateTokenPairCookies(c.Writer, tokenPair)
	if cuserr != nil {
		log.Printf("[AuthHandler] RefreshToken: Failed to generate new token pair: %s, RequestID: %s\n", cuserr.Message(), c.GetString("RequestID"))
		c.JSON(cuserr.Code(), dto.MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	log.Printf("[AuthHandler] RefreshToken: Token refreshed successfully for UserID: %d, Email: %s, RequestID: %s\n", userID, claims.Email, c.GetString("RequestID"))
	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Token berhasil diperbarui.",
	})
}

// @Summary User logout
// @Description Clear user authentication cookies
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.MessageResponse "Success response with message"
// @Router /auth/logout [post]
func (h *authHandler) Logout(c *gin.Context) {
	log.Printf("[AuthHandler] Logout: Started, RequestID: %s\n", c.GetString("RequestID"))
	h.authService.GenerateLogoutCookies(c.Writer)
	log.Printf("[AuthHandler] Logout: Cookies cleared, RequestID: %s\n", c.GetString("RequestID"))

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Logout berhasil.",
	})
}

// RegisterRoutes registers all auth routes
func (h *authHandler) RegisterRoutes(router *gin.RouterGroup) {
	log.Println("[AuthHandler] RegisterRoutes: Registering auth routes")
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/token/:type", h.RequestToken)
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
		authGroup.POST("/forget-password", h.ForgetPassword)
		authGroup.GET("/refresh-token", h.RefreshToken)
		authGroup.DELETE("/logout", h.Logout) // Changed to DELETE as per RESTful practices for logout
	}
	log.Println("[AuthHandler] RegisterRoutes: Auth routes registered")
}
