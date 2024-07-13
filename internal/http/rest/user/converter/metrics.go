package converter

import (
	"github.com/sashaaro/go-musthave-diploma-tpl/internal/domain/entity"
	updateInterface "github.com/sashaaro/go-musthave-diploma-tpl/internal/http/rest/update/interface"
)

func MetricsToMetricValueDTO(metrics entity.MetricJSON) updateInterface.GetMetricValueDto {
	return updateInterface.GetMetricValueDto{
		Type: metrics.MType,
		Name: metrics.ID,
	}
}
