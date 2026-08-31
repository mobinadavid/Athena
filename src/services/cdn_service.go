package services

import (
	"athena/src/api/errs"
	"athena/src/config"
	"athena/src/pkg/logger"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"sort"
	"strings"

	"github.com/go-resty/resty/v2"
)

type ICDNService interface {
	GetFile() interface{}
	GetFileWithTag(ctx context.Context, tag string, bucket string) (string, error)
	GetTagsString(tagMap map[string]string) string
	PutFile(ctx context.Context) (interface{}, error)
	UploadExcelToCDN(ctx context.Context) (interface{}, error)
	DeleteFile(ctx context.Context) error
	GetBytesByUrl(ctx context.Context, url string) ([]byte, error)
}

type CDNService struct {
	CDNUrl string
}

func NewCDNService() *CDNService {
	cdnUrl := config.GetInstance().Get("CDN_URL")
	return &CDNService{
		CDNUrl: cdnUrl,
	}
}

func (service *CDNService) GetFile() interface{} {
	return nil
}

func (service *CDNService) PutFile(ctx context.Context) (interface{}, error) {
	form := ctx.Value("form").(*multipart.Form)
	bucket := ctx.Value("bucket").(string)
	folder := ctx.Value("folder").(string)
	tags := ctx.Value("tags").(string)

	// fetch data
	url := fmt.Sprintf("%s/api/v1/buckets", config.GetInstance().Get("CDN_URL"))
	apikey := config.GetInstance().Get("CDN_API_KEY")

	// prepare client
	req := resty.New().R().
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("X-API-Key", apikey)

	// set form data
	req.SetFormData(map[string]string{
		"bucket": bucket,
		"folder": folder,
		"tag":    tags,
	})

	// Add the file to the request using SetFileReader
	for _, fileHeader := range form.File["files[]"] {
		file, err := fileHeader.Open()
		if err != nil {
			logger.LogErrorWithFieldsV2(ctx, "failed to open file", service, err)
			return nil, errs.CannotOpenFile
		}
		defer file.Close()

		// Add the file to the request using SetFileReader
		req = req.SetFileReader("files[]", fileHeader.Filename, file)
	}

	// send request
	resp, err := req.Post(url)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to send request using resty", service, err)
		return nil, errs.SomeThingWentWrong
	}
	return resp, nil
}

//func (service *CDNService) GetObjectsName(ctx context.Context) ([]string, error) {
//	bucket := ctx.Value("bucket").(string)
//
//	// fetch data
//	url := fmt.Sprintf("%s/api/v1/buckets/%s/objects", config.GetInstance().Get("CDN_URL"), bucket)
//	apikey := config.GetInstance().Get("CDN_API_KEY")
//
//	// prepare client
//	resp, err := resty.New().R().
//		SetHeaders(map[string]string{
//			"Content-Type": "application/json",
//			"Accept":       "application/json",
//			"X-API-Key":    apikey,
//		}).
//		Get(url)
//
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to send request using resty", service, err)
//		return nil, errs.SomeThingWentWrong
//	}
//
//	if !resp.IsSuccess() {
//		logger.LogErrorWithFieldsV2(ctx, "failed to send request using resty", service, nil)
//		return nil, errs.CDNServiceIsTemporaryDown
//	}
//
//	// Parse the successful response
//	var response dto.ListObjectsResponse
//	err = json.Unmarshal(([]byte(resp.String())), &response)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to unmarshal response", service, err)
//		return nil, errs.SomeThingWentWrong
//	}
//	// Extract object names
//	objectNames := []string{}
//	for _, object := range response.Data.Objects {
//		objectNames = append(objectNames, object.Info.UserTags.Name)
//	}
//
//	// Return the list of object names
//	return objectNames, nil
//}
//
//func (service *CDNService) UploadExcelToCDN(ctx context.Context) (interface{}, error) {
//	file := ctx.Value("file").(*excelize.File)
//	bucket := ctx.Value("bucket").(string)
//	folder := ctx.Value("folder").(string)
//	tags := ctx.Value("tags").(string)
//
//	// Convert the Excel file to a byte slice
//	var buf bytes.Buffer
//	if err := file.Write(&buf); err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to write Excel file to buffer", service, err)
//		return "", errs.SomeThingWentWrong
//	}
//
//	// Prepare the multipart form data to upload to CDN
//	// This will contain the actual file and its metadata
//	fileName := fmt.Sprintf("%s.xlsx", uuid.New().String()) // Generate a unique filename
//
//	// Create a new multipart form
//	var form bytes.Buffer
//	writer := multipart.NewWriter(&form)
//
//	// Add the metadata fields for CDN (if necessary)
//	writer.WriteField("bucket", bucket)
//	writer.WriteField("folder", folder)
//	writer.WriteField("tag", tags)
//
//	// Add the Excel file as part of the form
//	part, err := writer.CreateFormFile("files[]", fileName)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to create form file", service, err)
//		return "", errs.SomeThingWentWrong
//	}
//
//	// Write the Excel file content into the form
//	if _, err = part.Write(buf.Bytes()); err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to write file content to form", service, err)
//		return "", errs.SomeThingWentWrong
//	}
//
//	// Close the multipart writer
//	writer.Close()
//
//	// Set up the HTTP request to upload to CDN
//	url := fmt.Sprintf("%s/api/v1/buckets", config.GetInstance().Get("CDN_URL"))
//	apikey := config.GetInstance().Get("CDN_API_KEY")
//
//	// Prepare the request
//	req := resty.New().R().
//		SetHeader("Content-Type", writer.FormDataContentType()).SetHeader("X-API-Key", apikey).
//		SetBody(form.Bytes()) // Set the body as the multipart form data
//
//	// Send the request to upload the file
//	resp, err := req.Post(url)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to send request using resty", service, err)
//		return nil, errs.SomeThingWentWrong
//	}
//	return resp, nil
//}

//func (service *CDNService) GetFileWithTag(ctx context.Context, tag string, bucket string) (string, error) {
//	// fetch data
//	url := fmt.Sprintf("%s/api/v1/buckets/%s/tags/%s", config.GetInstance().Get("CDN_URL"), bucket, tag)
//	apikey := config.GetInstance().Get("CDN_API_KEY")
//
//	// prepare client
//	resp, err := resty.New().R().
//		SetHeaders(map[string]string{
//			"Content-Type": "application/json",
//			"Accept":       "application/json",
//			"X-API-Key":    apikey,
//		}).
//		Get(url)
//
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to send request using resty", service, err)
//		return "", errs.SomeThingWentWrong
//	}
//
//	if !resp.IsSuccess() {
//		logger.LogErrorWithFieldsV2(ctx, "failed to send request using resty", service, nil)
//		return "", errs.CDNServiceIsTemporaryDown
//	}
//
//	// Parse the successful response
//	var response dto.Response
//	err = json.Unmarshal(([]byte(resp.String())), &response)
//	if err != nil {
//		logger.LogErrorWithFieldsV2(ctx, "failed to unmarshal response", service, err)
//		return "", errs.SomeThingWentWrong
//	}
//
//	if len(response.Data.Objects) <= 0 {
//		return "", nil
//	}
//	return response.Data.Objects[0].Name, nil
//}

func (service *CDNService) DeleteFile(ctx context.Context) error {
	bucket := ctx.Value("bucket").(string)
	fileName := ctx.Value("file").(string)

	cdnURL := config.GetInstance().Get("CDN_URL")
	apiKey := config.GetInstance().Get("CDN_API_KEY")

	url := fmt.Sprintf("%s/api/v1/buckets/%s/files/%s", cdnURL, bucket, fileName)

	resp, err := resty.New().R().
		SetHeader("X-API-Key", apiKey).
		Delete(url)

	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to send request using resty", service, err)
		return errs.CDNServiceIsTemporaryDown
	}

	if resp.IsError() {
		logger.LogErrorWithFieldsV2(ctx, "CDN returned error", service, nil)
		return errs.SomeThingWentWrong
	}

	return nil
}

func (service *CDNService) GetTagsString(tagMap map[string]string) string {
	if len(tagMap) == 0 {
		return ""
	}

	keys := make([]string, 0, len(tagMap))
	for k := range tagMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for i, k := range keys {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(tagMap[k])
	}
	return b.String()
}

func (service *CDNService) GetBytesByUrl(ctx context.Context, url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to get file bytes from the url", service, err)
		return nil, errs.SomeThingWentWrong
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.LogErrorWithFieldsV2(ctx, "response status code not ok", service, nil)
		return nil, errs.SomeThingWentWrong
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.LogErrorWithFieldsV2(ctx, "failed to read response body", service, err)
		return nil, errs.SomeThingWentWrong
	}

	return data, nil
}

// GetImageAttachment generates a standard attachment response for an image URL using Pre-Signed URLs
//func GetImageAttachment(imageName string) models.Attachment {
//	if imageName != "" {
//		var cdnClient = cdn.NewCdnClient()
//		bucketName := config.GetInstance().Get("CDN_BUCKET_NAME")
//
//		preSignedData, err := cdnClient.GetPreSigned(bucketName, imageName)
//
//		if err == nil && preSignedData != nil && preSignedData.Data.URL != "" {
//			return models.Attachment{
//				Name: imageName,
//				URL:  base64.StdEncoding.EncodeToString([]byte(preSignedData.Data.URL)),
//			}
//		}
//	}
//
//	return models.Attachment{
//		Name: "",
//		URL:  "",
//	}
//}

// safeStr returns the dereferenced value of a string pointer, or "" if nil.
func SafeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
