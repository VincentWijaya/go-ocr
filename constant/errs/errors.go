package errs

import (
	"errors"
)

var (
	//errors
	Unauthorized             = errors.New("unauthorized")
	NoData                   = errors.New("cannot fetch requested data")
	BadRequest               = errors.New("bad Request")
	BadConfig                = errors.New("bad configuration on ocr")
	FaceNotFound             = errors.New("face is not registered in database")
	PlateNotRecognize        = errors.New("plate is not recognize")
	DriverFaceNotDetected    = errors.New("driver face is not detected")
	PlatenumberNotRegistered = errors.New("plate number is not registered")
)

var (
	//codes
	Success                           = "00"
	UnauthorizedErrorCode             = "03"
	NoDataCode                        = "-1"
	BadRequestCode                    = "04"
	BadConfigCode                     = "05"
	UndefinedErrorCode                = "99"
	FaceNotFoundErrorCode             = "40"
	PlateNotRecognizeErrorCode        = "41"
	DriverFaceNotDetectedErrorCode    = "42"
	PlatenumberNotRegisteredErrorCode = "43"

	//messages
	GeneralErrorMessage             = "Saat ini sedang terjadi gangguan pada system, silahkan coba beberapa saat lagi"
	NoDataMessage                   = "Data tidak di temukan"
	UnauthorizedMessage             = "Transaksi yang anda lakukan tidak sah"
	FaceNotFoundMessage             = "Wajah pengendara tidak terdaftar"
	PlateNotRecognizeMessage        = "Plat nomor tidak dapat terdeteksi, pastikan pengambilan gambar sudah benar"
	DriverFaceNotDetectedMessage    = "Wajah pengendara tidak dapat dideteksi, pastikan wajah pengendara terlihat jelas"
	PlatenumberNotRegisteredMessage = "Plat nomor tidak terdaftar disistem!"
)
