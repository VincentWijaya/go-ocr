package recog

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"

	"github.com/vincentwijaya/go-ocr/internal/app/domain"

	"github.com/Kagami/go-face"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vincentwijaya/go-ocr/pkg/log"
)

func GetFaceDescriptor(ctx context.Context, facePhotoLocation string) (descriptor [512]byte, err error) {
	logger := log.WithFields(log.Fields{"request_id": middleware.GetReqID(ctx)})

	rec, err := face.NewRecognizer("./testdata/models")
	if err != nil {
		logger.Errorf("Can't init face recognizer: %v", err)
		return
	}
	defer rec.Close()

	face, err := mustRecognizeSingleFile(rec, facePhotoLocation)
	if err != nil {
		logger.Errorf("Can't recognize face: %v", err)
		return
	}

	return descriptorToBytes(face.Descriptor), nil
}

func CompareFace(actualFaceDescriptor [512]byte, memberFaces []domain.Face, tolerance float32) uint {
	var (
		length     = len(memberFaces)
		categories = make([]int32, length)
		samples    = make([]face.Descriptor, length)
	)

	rec, err := face.NewRecognizer("./testdata/models")
	if err != nil {
		log.Errorf("Can't init face recognizer: %v", err)
		return 0
	}
	defer rec.Close()

	for i, f := range memberFaces {
		descriptor := bytesToDescriptor(f.Descriptor)
		samples[i] = descriptor
		categories[i] = int32(f.MemberID)
	}

	rec.SetSamples(samples, categories)

	var memberID = rec.ClassifyThreshold(bytesToDescriptor(actualFaceDescriptor), tolerance)

	return uint(memberID)
}

func descriptorToBytes(descriptor [128]float32) [512]byte {
	var result [512]byte

	var buffer = result[:0]

	for i := 0; i < 128; i++ {
		var bits uint32 = math.Float32bits(descriptor[i])

		buffer = append(
			buffer,
			byte(bits),
			byte(bits>>8),
			byte(bits>>16),
			byte(bits>>24),
		)
	}

	return result
}

func bytesToDescriptor(bytes [512]byte) [128]float32 {
	var result [128]float32

	var i = 0

	for j := 0; j < 512; j += 4 {
		result[i] = math.Float32frombits(
			uint32(bytes[j]) +
				uint32(bytes[j+1])<<8 +
				uint32(bytes[j+2])<<16 +
				uint32(bytes[j+3])<<24,
		)

		i += 1
	}

	return result
}

func mustRecognizeSingleFile(rec *face.Recognizer, filename string) (face.Face, error) {
	imageBytes, readFileErr := ioutil.ReadFile(filename)
	if readFileErr != nil {
		return face.Face{}, readFileErr
	}

	faces, recognizeErr := rec.Recognize(imageBytes)
	if recognizeErr != nil {
		return face.Face{}, recognizeErr
	}

	length := len(faces)
	if length != 1 {
		return face.Face{}, fmt.Errorf("expected 1 face on photo %s, got %d faces", filename, length)
	}

	return faces[0], nil
}
