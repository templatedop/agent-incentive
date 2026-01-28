//go:generate govalid ./marker.go

package test
import "time"

type Required struct {
	// +govalid:required
	Name string `validate:"required" json:"name"`

	// +govalid:required
	Age int `validate:"required" json:"age"`

	// +govalid:required
	Items []string `validate:"required" json:"items"`
}

type LT struct {
	// +govalid:lt=10
	Age int `validate:"lt=10" json:"age"`
}

type GT struct {
	// +govalid:gt=100
	Age int `validate:"gt=100" json:"age"`
}

type MaxLength struct {
	// +govalid:maxlength=50
	Name string `validate:"max=50" json:"name"`
}

type MinLength struct {
	// +govalid:minlength=3
	Name string `validate:"min=3" json:"name"`
}

type GTE struct {
	// +govalid:gte=18
	Age int `validate:"gte=18" json:"age"`
}

type LTE struct {
	// +govalid:lte=100
	Age int `validate:"lte=100" json:"age"`
}

type MaxItems struct {
	// +govalid:maxitems=5
	Items []string `validate:"max=5" json:"items"`

	// +govalid:maxitems=3
	Metadata map[string]string `json:"metadata"`

	// +govalid:maxitems=2
	ChanField chan int `json:"chan_field"`
}

type MinItems struct {
	// +govalid:minitems=2
	Items []string `validate:"min=2" json:"items"`

	// +govalid:minitems=1
	Metadata map[string]string `json:"metadata"`

	// +govalid:minitems=1
	ChanField chan int `json:"chan_field"`
}

type Alpha struct {
	// +govalid:alpha
	FirstName string `validate:"alpha" json:"first_name"`

	// +govalid:alpha
	LastName string `validate:"alpha" json:"last_name"`

	// +govalid:alpha
	CountryCode string `validate:"alpha" json:"country_code"`
}

// Custom types for enum testing
type UserRole string
type Priority int

type Enum struct {
	// String enum
	// +govalid:enum=admin,user,guest
	Role string `json:"role"`

	// Numeric enum
	// +govalid:enum=1,2,3
	Level int `json:"level"`

	// Custom string type enum
	// +govalid:enum=manager,developer,tester
	UserRole UserRole `json:"user_role"`

	// Custom numeric type enum
	// +govalid:enum=10,20,30
	Priority Priority `json:"priority"`
}

type Email struct {
	// +govalid:email
	Email string `validate:"email" json:"email"`
}

type UUID struct {
	// +govalid:uuid
	UUID string `validate:"uuid" json:"uuid"`
}

type URL struct {
	// +govalid:url
	URL string `validate:"url" json:"url"`
}

type CEL struct {
	// +govalid:cel=value >= 18
	Age int `json:"age"`

	// +govalid:cel=size(value) > 0
	Name string `json:"name"`

	// +govalid:cel=value > 0.0
	Score float64 `json:"score"`

	// +govalid:cel=value == true
	IsActive bool `json:"is_active"`
}

type CELCrossField struct {
	// Cross-field validation: Price must be less than MaxPrice
	// +govalid:cel=value < this.MaxPrice
	Price float64 `json:"price"`

	MaxPrice float64 `json:"max_price"`

	// Cross-field validation: Quantity * Price <= Budget
	// +govalid:cel=value * this.Price <= this.Budget
	Quantity float64 `json:"quantity"`

	Budget float64 `json:"budget"`
}

type Length struct {
	// +govalid:length=7
	Name string `validate:"len=7" json:"name"`
}

type Numeric struct {
	// +govalid:numeric
	Number string `json:"Number"`
}

type MultipleErrors struct {
	// +govalid:required
	URL string `validate:"required" json:"url"`

	// +govalid:maxlength=1
	TooLong string `validate:"max=1" json:"too_long"`
}

type IPV4 struct {
	// +govalid:ipv4
	IP string `validate:"ipv4" json:"ip"`
}

type IPV6 struct {
	// +govalid:ipv6
	IP string `validate:"ipv6" json:"ip"`
}

// New validators - Simple validators

type Min struct {
	// +govalid:min=10
	Age int `validate:"min=10" json:"age"`
}

type Eq struct {
	// +govalid:eq=active
	Status string `validate:"eq=active" json:"status"`

	// +govalid:eq=100
	Count int `validate:"eq=100" json:"count"`
}

type Ne struct {
	// +govalid:ne=admin
	Role string `validate:"ne=admin" json:"role"`

	// +govalid:ne=0
	Score int `validate:"ne=0" json:"score"`
}

type IsDefault struct {
	// +govalid:isdefault
	OptionalField string `validate:"isdefault" json:"optional_field"`

	// +govalid:isdefault
	OptionalNumber int `validate:"isdefault" json:"optional_number"`
}

type Boolean struct {
	// +govalid:boolean
	Flag string `validate:"boolean" json:"flag"`
}

type Lowercase struct {
	// +govalid:lowercase
	Username string `validate:"lowercase" json:"username"`
}

type OneOf struct {
	// +govalid:oneof=red green blue
	Color string `validate:"oneof=red green blue" json:"color"`

	// +govalid:oneof=1 2 3
	Level int `validate:"oneof=1 2 3" json:"level"`
}

type Number struct {
	// +govalid:number
	NumericString string `validate:"number" json:"numeric_string"`
}

type Alphanum struct {
	// +govalid:alphanum
	Code string `validate:"alphanum" json:"code"`
}

// String pattern validators

type ContainsAny struct {
	// +govalid:containsany=!@#$
	Password string `validate:"containsany=!@#$" json:"password"`
}

type Excludes struct {
	// +govalid:excludes=admin
	Username string `validate:"excludes=admin" json:"username"`
}

type ExcludesAll struct {
	// +govalid:excludesall=<>
	Comment string `validate:"excludesall=<>" json:"comment"`
}

type Unique struct {
	// +govalid:unique
	Tags []string `validate:"unique" json:"tags"`

	// +govalid:unique
	IDs []int `validate:"unique" json:"ids"`
}

// Format validators

type URI struct {
	// +govalid:uri
	Address string `validate:"uri" json:"address"`
}

type FQDN struct {
	// +govalid:fqdn
	Domain string `validate:"fqdn" json:"domain"`
}

type Latitude struct {
	// +govalid:latitude
	Lat string `validate:"latitude" json:"lat"`
}

type Longitude struct {
	// +govalid:longitude
	Lon string `validate:"longitude" json:"lon"`
}

type IsColour struct {
	// +govalid:iscolour
	Color string `validate:"iscolor" json:"color"`
}

// Duration validators

type MinDuration struct {
	// +govalid:minduration=1h
	Timeout time.Duration ` json:"timeout"`
}

type MaxDuration struct {
	// +govalid:maxduration=24h
	Interval time.Duration ` json:"interval"`
}

// Conditional required validators

type RequiredIf struct {
	Status string

	// +govalid:required_if=Status active
	ActiveField string `validate:"required_if=Status active" json:"active_field"`
}

type RequiredUnless struct {
	Status string

	// +govalid:required_unless=Status inactive
	ActiveField string `validate:"required_unless=Status inactive" json:"active_field"`
}

type RequiredWith struct {
	Email string

	// +govalid:required_with=Email
	EmailConfirmation string `validate:"required_with=Email" json:"email_confirmation"`
}

type RequiredWithAll struct {
	FirstName string
	LastName  string

	// +govalid:required_with_all=FirstName LastName
	FullName string `validate:"required_with_all=FirstName LastName" json:"full_name"`
}

type RequiredWithout struct {
	Phone string

	// +govalid:required_without=Phone
	Email string `validate:"required_without=Phone" json:"email"`
}

type RequiredWithoutAll struct {
	Phone string
	Fax   string

	// +govalid:required_without_all=Phone Fax
	Email string `validate:"required_without_all=Phone Fax" json:"email"`
}

// Conditional excluded validators

type ExcludedIf struct {
	Status string

	// +govalid:excluded_if=Status inactive
	InactiveField string `validate:"excluded_if=Status inactive" json:"inactive_field"`
}

type ExcludedUnless struct {
	Status string

	// +govalid:excluded_unless=Status active
	InactiveField string `validate:"excluded_unless=Status active" json:"inactive_field"`
}

type ExcludedWith struct {
	GuestMode string

	// +govalid:excluded_with=GuestMode
	AdminPanel string `validate:"excluded_with=GuestMode" json:"admin_panel"`
}

type ExcludedWithAll struct {
	ReadOnly  string
	Archived  string

	// +govalid:excluded_with_all=ReadOnly Archived
	EditButton string `validate:"excluded_with_all=ReadOnly Archived" json:"edit_button"`
}

type ExcludedWithout struct {
	Premium string

	// +govalid:excluded_without=Premium
	FreeFeature string `validate:"excluded_without=Premium" json:"free_feature"`
}

type ExcludedWithoutAll struct {
	FeatureA string
	FeatureB string

	// +govalid:excluded_without_all=FeatureA FeatureB
	ConflictingFeature string `validate:"excluded_without_all=FeatureA FeatureB" json:"conflicting_feature"`
}
