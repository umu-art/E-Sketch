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

const JWT_SIGNING_METHOD string = "HS256"
const JWT_DURATION_TIME = time.Hour * 48

const JWT_COOKIE_NAME string = "estu"

var SESSION_CHECK_EXCLUDED_PATH_SUFFIXES = []string{
	"/login",
	"/register",
	"/actuator",
}
