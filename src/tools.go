// EmikoTV ~ Written & maintained by Harvey S. Coombs
package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func SerializeToJSON[T any](object T) string {
	raw, err := json.Marshal(object)
	if err == nil {
		return string(raw)
	} else {
		return ""
	}
}

func GenerateSHA256(raw string) string {
	final := ""
	hash := sha256.New()

	hash.Write([]byte(raw))
	final = fmt.Sprintf("%x", hash.Sum(nil))

	return final
}

func EscapeHTML(raw string) string {
	escaped := raw

	escaped = strings.ReplaceAll(raw, ">", "&gt;")
	escaped = strings.ReplaceAll(raw, "<", "&lt;")
	escaped = strings.ReplaceAll(raw, "&", "&amp;")

	return escaped
}

func EscapeSQL(raw string) string {
	escaped := raw

	escaped = strings.ReplaceAll(raw, "'", "''")
	escaped = strings.ReplaceAll(raw, "\"", "\"\"")

	return escaped
}

func EscapeBoth(raw string) string {
	escaped := raw

	escaped = EscapeHTML(raw)
	escaped = EscapeSQL(raw)

	return escaped
}

func DetermineLocation(ip_addr string) string {
	country := ""
	cmd, err := exec.Command("whois", ip_addr).Output()

	if err != nil {
		log.Fatal(err)
	}

	result := string(cmd)
	startpos := (strings.Index(result, "country:") + 8)

	if startpos != -1 {
		result = result[startpos:len([]rune(result))]
		result = result[0:strings.Index(result, "\n")]

		country = strings.ToLower(strings.ReplaceAll(result, " ", ""))
	}

	return country
}

func FormatAsHTML(original string) string {
	original = strings.ReplaceAll(original, "\n", "<br/>")
	original = strings.ReplaceAll(original, "[c]", "<code>")
	original = strings.ReplaceAll(original, "[/c]", "</code>")
	original = strings.ReplaceAll(original, "[b]", "<b>")
	original = strings.ReplaceAll(original, "[/b]", "</b>")
	original = strings.ReplaceAll(original, "[i]", "</i>")

	return original
}

func StringToInt(raw string) int {
	integer := 0

	parsed, err := strconv.Atoi(raw)

	if err == nil {
		integer = parsed
	}

	return integer
}

func OpenFile(path string) string {
	file, err := os.ReadFile(path)

    if err != nil {
        return err.Error()
    }

	return string(file)
}

func OpenFileBytes(path string) []byte {
	file, err := os.ReadFile(path)

    if err != nil {
        return []byte{}
    }

	return file
}

func DirSize(path string) int {
	total := 0
	dir, err := os.ReadDir(path)

	if err == nil {
		for _, f := range dir {
			if !f.IsDir() {
				info, err := f.Info()

				if err == nil {
					total += int(info.Size())
				}
			}
		}
	}

	return total
}

func ListFilesInDir(path string) []string {
	files_list := []string{}

	files, err := os.ReadDir(path)

    if err == nil {
		for _, file := range files {
			files_list = append(files_list, file.Name())
		}
    }

	return files_list
}

func FileSize(path string) int {
	total := 0
	file, err := os.Stat(path)

	if err != nil {
		log.Fatal(err)
	} else {
		total = int(file.Size())
	}

	return total
}

func CreateFile(path string, content string) bool {
	result := false

	file, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	w, err := file.Write([]byte(content))

	if err == nil {
		result = (w > 0)
	}

	return result
}

func CreateFileFromBytes(path string, bytes []byte) bool {
	result := false

	file, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	w, err := file.Write(bytes)

	if err == nil {
		result = (w > 0)
	}

	return result
}

func DeleteFile(path string) bool {
	err := os.Remove(path)
	if err != nil {
		fmt.Printf("ERR! %s", err.Error())
	}
	
	return err == nil
}

func RenameFile(old string, new string) bool {
	DeleteFile(old)
	err := os.Rename(old, new)
	if err != nil {
		fmt.Printf("ERR! %s", err.Error())
	}

	return err == nil
}

func FileExists(path string) bool {
	_, err := os.OpenFile(path, os.O_RDWR, 0755)
	return !os.IsNotExist(err)
}

func CheckPlurality(amt int) string {
	if amt != 1 {
		return "s"
	} else {
		return ""
	}
}

func GenerateThumbnail(input_path string, output_path string) bool {
	cmdArgs := []string{
		"-i", input_path,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-q:v", "2",
		output_path,
	}

	cmd := exec.Command("ffmpeg", cmdArgs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	return true
}

func ConvertToMP4(input_path string, output_path string) bool {
	cmd := exec.Command("ffmpeg", "-i", input_path, output_path)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err == nil && FileExists(output_path)
}

func TrimVideo(target_file_path string, output_file_path string, duration float64) bool {
	cmd := exec.Command("ffmpeg",
		"-ss",
		"0", 
		"-i",  
		target_file_path,
		"-t",
		fmt.Sprint(duration),
		"-c:v",
		"libx264",
		"-c:a",
		"aac",
		output_file_path,
	)
	
	err := cmd.Run()
	return err == nil
}

func CompressVideo(target_file_path string, output_file_path string, scale string) bool {
	cmd := exec.Command("ffmpeg", "-i", target_file_path, "-vf", "scale=" + scale, "-c:v", "libx264", "-crf", "23", "-c:a", "aac", output_file_path)

	err := cmd.Run()
	return err == nil
}

func UploadFile(upload multipart.File, path string, filename string) string {
	defer upload.Close()

	DeleteFile(root + path + filename)
	stored_file, err := os.Create(root + path + filename)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer stored_file.Close()

	_, err = io.Copy(stored_file, upload)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return filename
}

func DetermineContentType(path string) string {
	if !FileExists(path) {
		return ""
	}

	ext := filepath.Ext(path)
	content_type := mime.TypeByExtension(ext)

	if len(content_type) == 0 {
		return "application/octet-stream"
	}

	return content_type
}