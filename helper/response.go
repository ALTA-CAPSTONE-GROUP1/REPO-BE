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
	currentLimit := limit
	if currentLimit <= 0 {
		currentLimit = totalData
	}

	currentOffset := offset
	if currentOffset <= 0 {
		currentOffset = 0
	}

	currentPage := 1
	if limit > 0 {
		currentPage = (offset / limit) + 1
	}

	totalPage := 1
	if limit > 0 {
		totalPage = int(math.Ceil(float64(totalData) / float64(limit)))
	}

	pagination := map[string]interface{}{
		"current_limit":  currentLimit,
		"current_offset": currentOffset,
		"current_page":   currentPage,
		"total_data":     totalData,
		"total_page":     totalPage,
	}

	return pagination
}
