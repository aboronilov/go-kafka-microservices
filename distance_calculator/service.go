package main

import (
	"math"

	"github.com/aboronilov/go-kafka-microservices/types"
)

type CalculatorServicer interface {
	CalculateDistance(types.OBUData) (float64, error)
}

type CalculatorService struct {
	points [][]float64
}

func NewCalculatorService() CalculatorServicer {
	return &CalculatorService{
		points: make([][]float64, 0),
	}
}

func (s *CalculatorService) CalculateDistance(data types.OBUData) (float64, error) {
	distance := 0.0
	if len(s.points) > 0 {
		prevPoints := s.points[len(s.points)-1]
		distance = s.getDistance(prevPoints[0], prevPoints[1], data.Lat, data.Long)
	}
	s.points = append(s.points, []float64{data.Lat, data.Long})

	return distance, nil
}

func (s *CalculatorService) getDistance(x1, x2, y1, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
