package a

type Markers struct {
	Required       string `validate:"required"`       // want Required:`Identifier: "govalid:required", Expressions: {no expressions}`
	Min            int    `validate:"lt=10"`          // want Min:`Identifier: "govalid:lt", Expressions: {govalid:lt: 10}`
	RequiredAndMin string `validate:"required,lt=10"` // want RequiredAndMin:`Identifier: "govalid:required", Expressions: {no expressions}` // want RequiredAndMin:`Identifier: "govalid:lt", Expressions: {govalid:lt: 10}`
}

type TypeLevelMarkers struct {
	String         string
	RequiredString int    `validate:"required"` // want RequiredString:`Identifier: "govalid:required", Expressions: {no expressions}`
	Email          string `validate:"email"`    // want Email:`Identifier: "govalid:email", Expressions: {no expressions}`
}
