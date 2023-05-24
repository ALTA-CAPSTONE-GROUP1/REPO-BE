package helper

import "math"

func ResponseFormat(code int, msg string, data any) (int, map[string]any) {
	res := map[string]any{}
	res["code"] = code
	res["message"] = msg

	if data != nil {
		res["data"] = data
	}

	return code, res
}

func ReponseFormatWithMeta(code int, msg string, data any, meta any) (int, map[string]any) {
	res := map[string]any{}
	res["code"] = code
	res["message"] = msg

	if meta != nil {
		res["meta"] = meta
	}
	if data != nil {
		res["data"] = data
	}

	return code, res
}

func Pagination(limit, offset, totalData int) map[string]interface{} {

	if offset < totalData {
		endIndex := offset + limit
		if endIndex > totalData {
			endIndex = totalData
		}
	}

	totalPage := 1
	if totalData > 0 {
		totalPage = int(math.Ceil(float64(totalData) / float64(limit)))
	}
	currentPage := int(math.Ceil(float64(offset+1) / float64(limit)))

	if currentPage > totalPage {
		currentPage = totalPage
	}

	pagination := map[string]interface{}{
		"current_limit":  limit,
		"current_offset": offset,
		"current_page":   currentPage,
		"total_data":     totalData,
		"total_page":     totalPage,
	}

	return pagination
}
