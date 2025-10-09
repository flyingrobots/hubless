package mock

type layoutProfile struct {
	id            string
	name          string
	twoColumn     bool
	listRatio     int
	previewRatio  int
	boardMinWidth int
}

// Name returns the human-readable label for the profile.
func (p layoutProfile) Name() string {
	if p.name == "" {
		return p.id
	}
	return p.name
}

func profileForWidth(width int) layoutProfile {
	switch {
	case width <= 0:
		return layoutProfile{
			id:            "md",
			name:          "medium",
			twoColumn:     true,
			listRatio:     5,
			previewRatio:  7,
			boardMinWidth: 28,
		}
	case width < 100:
		return layoutProfile{
			id:            "sm",
			name:          "small",
			twoColumn:     false,
			listRatio:     1,
			previewRatio:  1,
			boardMinWidth: 24,
		}
	case width < 140:
		return layoutProfile{
			id:            "md",
			name:          "medium",
			twoColumn:     true,
			listRatio:     5,
			previewRatio:  7,
			boardMinWidth: 28,
		}
	default:
		return layoutProfile{
			id:            "lg",
			name:          "large",
			twoColumn:     true,
			listRatio:     4,
			previewRatio:  7,
			boardMinWidth: 32,
		}
	}
}
