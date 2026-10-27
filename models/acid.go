package models

import "fmt"

// MinioBaseURL - базовый URL бакета acid-media в MinIO.
const MinioBaseURL = "http://localhost:9000/acid-media/"

// Acid - услуга предметной области:
type Acid struct {
	ID                    int
	ReactionName          string   // "CaCO3 + 2HCl"
	AcidName              string   // "Соляная кислота"
	AcidFormula           string   // "HCl"
	MolarMass             float64  // молярная масса кислоты, г/моль
	Concentration         float64  // концентрация кислоты, %
	CalciumCarbonateMass  float64  // масса CaCO3, г
	Description           string
	Likes                 []int    // ID пользователей, поставивших лайк
	Status                string   // draft | published | deleted
	ImageKey              string
	VideoKey              string
}

// ImageURL - прямая ссылка на изображение в MinIO.
func (a Acid) ImageURL() string {
	if a.ImageKey == "" {
		return ""
	}
	return MinioBaseURL + a.ImageKey
}

// VideoURL - прямая ссылка на видео в MinIO.
func (a Acid) VideoURL() string {
	if a.VideoKey == "" {
		return ""
	}
	return MinioBaseURL + a.VideoKey
}

// LikeCount - количество лайков (по коллекции ID пользователей).
func (a Acid) LikeCount() int {
	return len(a.Likes)
}

// FormatMolarMass - отображение молярной массы.
func FormatMolarMass(m float64) string {
	return fmt.Sprintf("%.2f г/моль", m)
}

// Acids - коллекция услуг (5 опубликованных + 1 черновик).
var Acids = []Acid{
	{
		ID:                   1,
		ReactionName:         "CaCO3 + 2HCl",
		AcidName:             "Соляная кислота",
		AcidFormula:          "HCl",
		MolarMass:            36.46,
		Concentration:        37,
		CalciumCarbonateMass: 50,
		Description:          "Выделение CO2 при взаимодействии карбоната кальция с соляной кислотой. Уравнение: CaCO3 + 2HCl → CaCl2 + CO2↑ + H2O.",
		Likes:                []int{101, 102, 103, 104, 105, 106, 107, 108, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124},
		Status:               "published",
		ImageKey:             "hcl_acid.jpg",
		VideoKey:             "hcl_reaction.mp4",
	},
	{
		ID:                   2,
		ReactionName:         "CaCO3 + H2SO4",
		AcidName:             "Серная кислота",
		AcidFormula:          "H2SO4",
		MolarMass:            98.08,
		Concentration:        98,
		CalciumCarbonateMass: 100,
		Description:          "Выделение CO2 при взаимодействии карбоната кальция с серной кислотой. Уравнение: CaCO3 + H2SO4 → CaSO4 + CO2↑ + H2O.",
		Likes:                []int{201, 202, 203, 204, 205, 206, 207, 208, 217, 218, 219, 220, 221, 222, 223, 224, 225},
		Status:               "published",
		ImageKey:             "h2so4_acid.jpg",
		VideoKey:             "h2so4_reaction.mp4",
	},
	{
		ID:                   3,
		ReactionName:         "CaCO3 + 2CH3COOH",
		AcidName:             "Уксусная кислота",
		AcidFormula:          "CH3COOH",
		MolarMass:            60.05,
		Concentration:        70,
		CalciumCarbonateMass: 25,
		Description:          "Выделение CO2 при взаимодействии карбоната кальция с уксусной кислотой. Уравнение: CaCO3 + 2CH3COOH → Ca(CH3COO)2 + CO2↑ + H2O.",
		Likes:                []int{301, 302, 303, 304, 305, 306, 307, 308, 309, 310, 311, 312, 313, 314, 315},
		Status:               "published",
		ImageKey:             "ch3cooh_acid.jpg",
		VideoKey:             "ch3cooh_reaction.mp4",
	},
	{
		ID:                   4,
		ReactionName:         "3CaCO3 + 2H3PO4",
		AcidName:             "Фосфорная кислота",
		AcidFormula:          "H3PO4",
		MolarMass:            98.00,
		Concentration:        85,
		CalciumCarbonateMass: 150,
		Description:          "Выделение CO2 при взаимодействии карбоната кальция с фосфорной кислотой. Уравнение: 3CaCO3 + 2H3PO4 → Ca3(PO4)2 + 3CO2↑ + 3H2O.",
		Likes:                []int{401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 418, 419, 420, 421, 422, 423, 424, 425, 426, 427, 428, 429, 430, 431, 432, 433, 434, 435, 436, 437, 438, 439, 440, 441, 442},
		Status:               "published",
		ImageKey:             "h3po4_acid.jpg",
		VideoKey:             "h3po4_reaction.mp4",
	},
	{
		ID:                   5,
		ReactionName:         "CaCO3 + 2HNO3",
		AcidName:             "Азотная кислота",
		AcidFormula:          "HNO3",
		MolarMass:            63.01,
		Concentration:        68,
		CalciumCarbonateMass: 0,
		Description:          "",
		Likes:                []int{},
		Status:               "draft",
		ImageKey:             "",
		VideoKey:             "",
	},
}

// PublishedAcids - только опубликованные.
func PublishedAcids() []Acid {
	out := make([]Acid, 0, len(Acids))
	for _, a := range Acids {
		if a.Status == "published" {
			out = append(out, a)
		}
	}
	return out
}

// DraftAcid - единственный черновик.
func DraftAcid() Acid {
	for _, a := range Acids {
		if a.Status == "draft" {
			return a
		}
	}
	return Acid{}
}

// FindByID - поиск опубликованной услуги по ID.
func FindByID(id int) (Acid, bool) {
	for _, a := range Acids {
		if a.ID == id && a.Status == "published" {
			return a, true
		}
	}
	return Acid{}, false
}

// NextPublished - следующая опубликованная после указанного ID.
// Если дошли до конца - циклически возвращаемся к первой.
func NextPublished(id int) Acid {
	pub := PublishedAcids()
	if len(pub) == 0 {
		return Acid{}
	}
	for i, a := range pub {
		if a.ID == id {
			return pub[(i+1)%len(pub)]
		}
	}
	return pub[0]
}

// FilterByMolarMassRange - фильтрация по молярной массе в диапазоне [min, max].
// Если min <= 0 и max <= 0 - возвращаем все опубликованные.
func FilterByMolarMassRange(min, max float64) []Acid {
	pub := PublishedAcids()
	if min <= 0 && max <= 0 {
		return pub
	}
	out := make([]Acid, 0, len(pub))
	for _, a := range pub {
		if a.MolarMass >= min && (max <= 0 || a.MolarMass <= max) {
			out = append(out, a)
		}
	}
	return out
}

// FirstPublishedID - ID первой опубликованной услуги (для редиректа с /feed без параметров).
func FirstPublishedID() int {
	pub := PublishedAcids()
	if len(pub) == 0 {
		return 0
	}
	return pub[0].ID
}