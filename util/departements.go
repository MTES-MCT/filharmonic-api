package util

import (
	"strings"
)

/*
Renvoie le code INSEE du département en fonction du code INSEE d'une commune
*/
func GetCodeDepartementFromCodeInseeCommune(codePostal string) string {
	if codePostal != "" {
		if strings.HasPrefix(codePostal, "97") { // cas outre-mer
			return codePostal[:3]
		} else {
			return codePostal[:2]
		}
	}
	return ""
}
