package endpoint

import (
	"captcha-backend/config"
	"captcha-backend/utils"
	"crypto/sha1"
	"encoding/hex"
	"github.com/go-redis/redis"
	"github.com/valyala/fasthttp"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"time"
)

type ICaptchaService interface {
	CreateCaptcha(word string, width int, height int) (*image.RGBA, error)
}

type Handler struct {
	redis          *redis.Client
	captchaService ICaptchaService
	config         *config.Config
}

func NewHandler(redis *redis.Client, captchaService ICaptchaService, config *config.Config) *Handler {
	return &Handler{
		redis:          redis,
		captchaService: captchaService,
		config:         config,
	}
}

const ttlSessionTime = 120 * time.Minute

func (h *Handler) GenerateSession(ctx *fasthttp.RequestCtx) {
	word := h.config.Words[rand.Intn(len(h.config.Words))]
	captcha, err := h.captchaService.CreateCaptcha(word, 500, 200)
	if err != nil {
		SendErrorResponse(ctx, fasthttp.StatusInternalServerError, []byte("error generating captcha "+err.Error()))
		return
	}

	userSessionKey := utils.RandStringRunes(16)

	expTime := time.Now().Unix()

	cmd := h.redis.Set(userSessionKey+"-sessionKey", 0, ttlSessionTime)
	if cmd.Err() != nil {
		SendErrorResponse(ctx, fasthttp.StatusInternalServerError, []byte("error setting session key in redis "+cmd.Err().Error()))
		return
	}

	hasher := sha1.New()
	hasher.Write([]byte(word))
	sha := hasher.Sum(nil)

	cypherText := word + userSessionKey

	for len([]byte(cypherText))%16 != 0 {
		cypherText += "0"
	}

	sha128 := sha[:16]

	captchaKey, err := utils.EncryptAES([]byte(hex.EncodeToString(sha128)), cypherText)
	if err != nil {
		SendErrorResponse(ctx, fasthttp.StatusInternalServerError, []byte("error generating aes session key "+err.Error()))
		return
	}

	outFile, err := os.Create("endpoint/captcha_images/" + captchaKey + ".jpeg")
	if err != nil {
		SendErrorResponse(ctx, fasthttp.StatusInternalServerError, []byte("error creating image "+err.Error()))
		return
	}

	err = jpeg.Encode(outFile, captcha, nil)
	if err != nil {
		SendErrorResponse(ctx, fasthttp.StatusInternalServerError, []byte("error encoding captcha to jpeg "+err.Error()))
		return
	}

	SendSuccessResponse(ctx, struct {
		Sha128         string
		UserSessionKey string
		CaptchaKey     string `json:"captcha_key"`
		ExpTime        int64  `json:"exp_time"`
	}{
		Sha128:         hex.EncodeToString(sha128),
		UserSessionKey: userSessionKey,
		CaptchaKey:     captchaKey,
		ExpTime:        expTime,
	})
}
