// Description: util
// Author: ZHU HAIHUA
// Since: 2016-06-17 11:42
package filepath

import "path/filepath"

// Extract the file name without extension.
// eg: path/to/file.txt -> file
// eg: path/to/file -> file
func BaseName(filePath string) string {
	if filePath == "" {
		return ""
	}
	ext := filepath.Ext(filePath)
	base := filepath.Base(filePath)
	return base[:len(base) - len(ext)]
}
