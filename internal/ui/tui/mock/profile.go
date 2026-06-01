package mock

// layoutProfile encapsulates responsive settings derived from Stickers breakpoints.
type layoutProfile struct {
	id   string
	name string
}

// Name returns the human-readable label for the profile.
func (p layoutProfile) Name() string {
	if p.name == "" {
		return p.id
	}
	return p.name
}

// profileForWidth returns the responsive layoutProfile for the given width.
// It maps widths < 100 to the "sm" (small) profile, widths >= 100 and < 140 to
// the "md" (medium) profile, and widths >= 140 to the "lg" (large) profile.
func profileForWidth(width int) layoutProfile {
	switch {
	case width < 100:
		return layoutProfile{id: "sm", name: "small"}
	case width < 140:
		return layoutProfile{id: "md", name: "medium"}
	default:
		return layoutProfile{id: "lg", name: "large"}
	}
}
