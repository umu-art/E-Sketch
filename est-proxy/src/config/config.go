package config

import (
	"os"
	"time"
)

var EST_BACK_URL = os.Getenv("EST_BACK_URL")
var EST_PREVIEW_URL = os.Getenv("EST_PREVIEW_URL")

var POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
var POSTGRES_USERNAME = os.Getenv("POSTGRES_USERNAME")
var POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
var POSTGRES_PORT = os.Getenv("POSTGRES_PORT")
var POSTGRES_DATABASE = os.Getenv("POSTGRES_DATABASE")

var RABBITMQ_USERNAME = os.Getenv("RABBITMQ_USERNAME")
var RABBITMQ_PASSWORD = os.Getenv("RABBITMQ_PASSWORD")
var RABBITMQ_HOST = os.Getenv("RABBITMQ_HOST")
var RABBITMQ_PORT = os.Getenv("RABBITMQ_PORT")
var RABBITMQ_FIGURE_TOPIC_EXCHANGE = "figure_change"
var RABBITMQ_MARKER_TOPIC_EXCHANGE = "marker_change"

var GPT_API_PATH = os.Getenv("GPT_API_PATH")
var GPT_API_TOKEN = os.Getenv("GPT_API_TOKEN")

var JWT_SECRET = os.Getenv("JWT_SECRET")

var JWT_SIGNING_METHOD string = "HS256"
var JWT_DURATION_TIME = time.Hour * 48

var JWT_COOKIE_NAME string = "estu"

var BUFFERED_FIGURE_EXPIRATION_TIME = 500 * time.Millisecond

var SMTP_SERVER = os.Getenv("SMTP_SERVER")
var SMTP_PORT = os.Getenv("SMTP_PORT")
var SMTP_EMAIL = os.Getenv("SMTP_EMAIL")
var SMTP_PASSWORD = os.Getenv("SMTP_PASSWORD")
var SMTP_DKIM_KEY_FILE = os.Getenv("SMTP_DKIM_KEY_FILE")

var REDIS_URL = os.Getenv("REDIS_URL")
var REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
var REDIS_DB = os.Getenv("REDIS_DB")

var REDIS_EXPIRATION_TIME = 5 * time.Minute

var CONFIRM_URL string = "https://e-sketch.ru/auth/confirm"

var SESSION_CHECK_EXCLUDED_PATH_SUFFIXES = []string{
	"/login",
	"/register",
	"/confirm",
	"/actuator",
}
