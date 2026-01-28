# Supported Markers

govalid supports the following markers:

## `govalid:required`
- **Description**: Ensures that the field is not empty or nil.
- **Example**:
  ```go
  // +govalid:required
  type User struct {
      Username string `json:"username"`
  }
  ```
- **Generated Code**:
  ```go
  func ValidateUser(t *User) error {
      if t == nil {
          return ErrNilUser
      }

      if t.Username == "" {
          return ErrUsernameRequiredValidation
      }

      return nil
  }
  ```

## `govalid:lt`
- **Description**: Ensures that a numeric field is less than a specified value.
- **Example**:
  ```go
  // +govalid:lt=18
  type Profile struct {
      Age int `json:"age"`
  }
  ```
- **Generated Code**:
  ```go
  func ValidateProfile(t *Profile) error {
      if t == nil {
          return ErrNilProfile
      }

      if t.Age >= 18 {
          return ErrAgeLtValidation
      }

      return nil
  }
  ```

## `govalid:lte`
- **Description**: Ensures that a numeric field is less than or equal to a specified value.
- **Example**:
  ```go
  // +govalid:lte=65
  type Profile struct {
      Age int `json:"age"`
  }
  ```
- **Generated Code**:
  ```go
  func ValidateProfile(t *Profile) error {
      if t == nil {
          return ErrNilProfile
      }

      if !(t.Age <= 65) {
          return ErrAgeLTEValidation
      }

      return nil
  }
  ```

## `govalid:gt`
- **Description**: Ensures that a numeric field is greater than a specified value.
- **Example**:
  ```go
  // +govalid:gt=100
  type Profile struct {
      Age int `json:"age"`
  }
  ```
- **Generated Code**:
  ```go
  func ValidateProfile(t *Profile) error {
      if t == nil {
          return ErrNilProfile
      }

      if t.Age <= 100 {
          return ErrAgeGtValidation
      }

      return nil
  }
  ```

## `govalid:gte`
- **Description**: Ensures that a numeric field is greater than or equal to a specified value.
- **Example**:
  ```go
  // +govalid:gte=18
  type Profile struct {
      Age int `json:"age"`
  }
  ```
- **Generated Code**:
  ```go
  func ValidateProfile(t *Profile) error {
      if t == nil {
          return ErrNilProfile
      }

      if !(t.Age >= 18) {
          return ErrAgeGTEValidation
      }

      return nil
  }
  ```

## `govalid:maxlength`
- **Description**: Ensures that a string field's length does not exceed the specified maximum value (Unicode-aware).
- **Example**:
  ```go
  type User struct {
      // +govalid:maxlength=50
      Username string `json:"username"`
  }
  ```
- **Generated Code**:
  ```go
  func ValidateUser(t *User) error {
      if t == nil {
          return ErrNilUser
      }

      if utf8.RuneCountInString(t.Username) > 50 {
          return ErrUsernameMaxLengthValidation
      }

      return nil
  }
  ```

## `govalid:minlength`
- **Description**: Ensures that a string field's length is at least the specified minimum value (Unicode-aware).
- **Example**:
  ```go
  type User struct {
      // +govalid:minlength=3
      Username string `json:"username"`
  }
  ```
- **Generated Code**:
  ```go
  func ValidateUser(t *User) error {
      if t == nil {
          return ErrNilUser
      }

      if utf8.RuneCountInString(t.Username) < 3 {
          return ErrUsernameMinLengthValidation
      }

      return nil
  }
  ```

## `govalid:length`
- **Description**: Ensures that a string field has exactly the specified length (Unicode-aware).
- **Example**:
  ```go
  type User struct {
      // +govalid:length=7
      Name string `json:"name"`
  }
  ```
- **Generated Code**:
  ```go
  func ValidateUser(t *User) error {
      if t == nil {
          return ErrNilUser
      }

      if utf8.RuneCountInString(t.Name) != 7 {
          return ErrUserNameLengthValidation
      }

      return nil
  }
  ```
- **Note**: Uses `utf8.RuneCountInString()` for proper Unicode character counting, ensuring accurate validation for international characters and emojis.

## `govalid:maxitems`
- **Description**: Ensures that a collection field's length does not exceed the specified maximum number of items. Supports slice, array, map, and channel types.
- **Example**:
  ```go
  type Collection struct {
      // +govalid:maxitems=10
      Items []string `json:"items"`
      
      // +govalid:maxitems=5
      Metadata map[string]string `json:"metadata"`
  }
  ```
- **Generated Code**:
  ```go
  func ValidateCollection(t *Collection) error {
      if t == nil {
          return ErrNilCollection
      }

      if len(t.Items) > 10 {
          return ErrItemsMaxItemsValidation
      }

      if len(t.Metadata) > 5 {
          return ErrMetadataMaxItemsValidation
      }

      return nil
  }
  ```

## `govalid:minitems`
- **Description**: Ensures that a collection field's length is at least the specified minimum number of items. Supports slice, array, map, and channel types.
- **Example**:
  ```go
  type Collection struct {
      // +govalid:minitems=1
      Items []string `json:"items"`
      
      // +govalid:minitems=2
      Tags []string `json:"tags"`
  }
  ```
- **Generated Code**:
  ```go
  func ValidateCollection(t *Collection) error {
      if t == nil {
          return ErrNilCollection
      }

      if len(t.Items) < 1 {
          return ErrItemsMinItemsValidation
      }

      if len(t.Tags) < 2 {
          return ErrTagsMinItemsValidation
      }

      return nil
  }
  ```

## `govalid:enum`
- **Description**: Ensures that a field value is within a specified set of allowed values. Supports string, numeric, and custom types with comparable values. Values should be comma-separated.
- **Example**:
  ```go
  // Custom types
  type UserRole string
  type Priority int

  type User struct {
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
  ```
- **Generated Code**:
  ```go
  func ValidateUser(t *User) error {
      if t == nil {
          return ErrNilUser
      }

      if t.Role != "admin" && t.Role != "user" && t.Role != "guest" {
          return ErrRoleEnumValidation
      }

      if t.Level != 1 && t.Level != 2 && t.Level != 3 {
          return ErrLevelEnumValidation
      }

      if t.UserRole != "manager" && t.UserRole != "developer" && t.UserRole != "tester" {
          return ErrUserRoleEnumValidation
      }

      if t.Priority != 10 && t.Priority != 20 && t.Priority != 30 {
          return ErrPriorityEnumValidation
      }

      return nil
  }
  ```

## `govalid:email`
- **Description**: Ensures that a string field is a valid email address using HTML5-compliant validation.
- **Example**:
  ```go
  type User struct {
      // +govalid:email
      Email string `validate:"email" json:"email"`
  }
  ```
- **Generated Code**:
  ```go
  func ValidateUser(t *User) error {
      if t == nil {
          return ErrNilUser
      }

      if !emailRegex.MatchString(t.Email) {
          return ErrEmailEmailValidation
      }

      return nil
  }
  ```

## `govalid:url`
- **Description**: Ensures that a string field is a valid URL using HTTP/HTTPS protocol validation.
- **Example**:
  ```go
  type Resource struct {
      // +govalid:url
      URL string `validate:"url" json:"url"`
  }
  ```
- **Generated Code**:
  ```go
  func ValidateResource(t *Resource) error {
      if t == nil {
          return ErrNilResource
      }

      if !validationhelper.IsValidURL(t.URL) {
          return ErrURLURLValidation
      }

      return nil
  }
  ```

## `govalid:uuid`
- **Description**: Ensures that a string field is a valid UUID following RFC 4122 format.
- **Example**:
  ```go
  type Resource struct {
      // +govalid:uuid
      ID string `json:"id"`
  }
  ```
- **Generated Code**:
  ```go
  func ValidateResource(t *Resource) error {
      if t == nil {
          return ErrNilResource
      }

      if !isValidUUID(t.ID) {
          return ErrIDUUIDValidation
      }

      return nil
  }
  ```

## `govalid:cel`
- **Description**: Validates fields using Google's Common Expression Language (CEL) for complex validation logic.
- **Available Variables**: 
  - `value`: The current field value being validated
- **Example**:
  ```go
  type Config struct {
      // Numeric validation
      // +govalid:cel=value > 0.0 && value <= 100.0
      Score float64 `json:"score"`
      
      // String length validation
      // +govalid:cel=size(value) >= 3 && size(value) <= 50
      Username string `json:"username"`
      
      // Pattern matching
      // +govalid:cel=value.contains('@')
      Email string `json:"email"`
      
      // List validation
      // +govalid:cel=size(value) > 0 && size(value) <= 10
      Tags []string `json:"tags"`
  }
  ```
- **Generated Code**:
  ```go
  func ValidateConfig(t *Config) error {
      if t == nil {
          return ErrNilConfig
      }

      if !validationhelper.IsValidCEL("value > 0.0 && value <= 100.0", t.Score, t) {
          return ErrScoreCELValidation
      }

      if !validationhelper.IsValidCEL("size(value) >= 3 && size(value) <= 50", t.Username, t) {
          return ErrUsernameCELValidation
      }

      if !validationhelper.IsValidCEL("value.contains('@')", t.Email, t) {
          return ErrEmailCELValidation
      }

      if !validationhelper.IsValidCEL("size(value) > 0 && size(value) <= 10", t.Tags, t) {
          return ErrTagsCELValidation
      }

      return nil
  }
  ```
- **Note**: CEL validation follows govalid's zero-reflection philosophy. Cross-field validation (accessing other struct fields) is not supported.

## `govalid:alpha`
- **Description**: Ensures that a string field is alphabetical, i.e. all its characters belong to the english alphabet.
- **Example**:
  ```go
    type User struct {
        // +govalid:alpha
        FirstName string `json:"first_name"`
    }
    ```
- **Generated Code**:
  ```go
  func ValidateUser(t *User) error {
    if t == nil {
        return ErrNilUser
    }

    if !validationhelper.IsValidAlpha(t.FirstName) {
        return ErrUserFirstNameAlphaValidation
    }

    return nil
  }
  ```

## `govalid:numeric`

- **Description**: Ensures that a string field contains only digit characters (`0`–`9`). Leading zeros are allowed. Negative signs, decimal points, or scientific notation are **not** accepted.

- **Example**:

  ```go
  type Payload struct {
      // +govalid:numeric
      Phone string `json:"phone"`
  }
  ```

- **Generated Code**:

  ```go
  func ValidatePayload(t *Payload) error {
      if t == nil {
          return ErrNilPayload
      }

      if !validationhelper.IsNumeric(t.Phone) {
          return ErrPayloadPhoneNumericValidation
      }

      return nil
  }
  ```

## `govalid:ipv4`

- **Description**: Ensure that a string field is a valid RFC 791-compliant IPv4 address.

- **Example**:

  ```go
  type Request struct {
      // +govalid:ipv4
      IP string `json:"ip"`
  }
  ```

- **Generated Code**:

  ```go
  func ValidateRequest(t *Request) error {
      if t == nil {
          return ErrNilRequest
      }

      var errs govaliderrors.ValidationErrors

      if ip := net.ParseIP(t.IP); ip == nil || ip.To4() == nil {
          err := ErrRequestIPIpv4Validation
          err.Value = t.IP
          errs = append(errs, err)
      }

      if len(errs) > 0 {
          return errs
      }
      return nil
  }

  ```

## `govalid:ipv6`

- **Description**: Ensure that a string field is a valid RFC 4291-compliant IPv6 address.

- **Example**:

  ```go
  type Request struct {
      // +govalid:ipv6
      IP string `json:"ip"`
  }
  ```

- **Generated Code**:

  ```go
  func ValidateRequest(t *Request) error {
      if t == nil {
          return ErrNilRequest
      }

      var errs govaliderrors.ValidationErrors

      if ip := net.ParseIP(t.IP); ip == nil || ip.To4() != nil {
          err := ErrRequestIPIpv6Validation
          err.Value = t.IP
          errs = append(errs, err)
      }

      if len(errs) > 0 {
          return errs
      }
      return nil
  }

  ```
## `govalid:min`
- **Description**: Ensures that a numeric field is greater than or equal to a specified minimum value.
- **Example**:
  ```go
  type Product struct {
      // +govalid:min=10
      Price int `json:"price"`
  }
  ```
- **Generated Code**:
  ```go
  if !(t.Price >= 10) {
      return ErrPriceMinValidation
  }
  ```

## `govalid:eq`
- **Description**: Ensures that a field equals a specific value. Supports strings and numeric types.
- **Example**:
  ```go
  type Status struct {
      // +govalid:eq=active
      State string `json:"state"`
      // +govalid:eq=100
      Count int `json:"count"`
  }
  ```
- **Generated Code**:
  ```go
  if !(t.State == "active") {
      return ErrStateEqValidation
  }
  if !(t.Count == 100) {
      return ErrCountEqValidation
  }
  ```

## `govalid:ne`
- **Description**: Ensures that a field does not equal a specific value. Supports strings and numeric types.
- **Example**:
  ```go
  type User struct {
      // +govalid:ne=admin
      Role string `json:"role"`
      // +govalid:ne=0
      Score int `json:"score"`
  }
  ```
- **Generated Code**:
  ```go
  if !(t.Role != "admin") {
      return ErrRoleNeValidation
  }
  if !(t.Score != 0) {
      return ErrScoreNeValidation
  }
  ```

## `govalid:isdefault`
- **Description**: Ensures that a field has its zero/default value (opposite of required).
- **Example**:
  ```go
  type Optional struct {
      // +govalid:isdefault
      OptionalField string `json:"optional_field"`
  }
  ```
- **Generated Code**:
  ```go
  if t.OptionalField != "" {
      return ErrOptionalFieldIsDefaultValidation
  }
  ```

## `govalid:boolean`
- **Description**: Ensures that a string field represents a valid boolean value. Accepts: true, false, 1, 0, yes, no, on, off (case-insensitive).
- **Example**:
  ```go
  type Settings struct {
      // +govalid:boolean
      Enabled string `json:"enabled"`
  }
  ```
- **Generated Code**:
  ```go
  if !validationhelper.IsValidBoolean(t.Enabled) {
      return ErrEnabledBooleanValidation
  }
  ```

## `govalid:lowercase`
- **Description**: Ensures that a string field contains only lowercase characters.
- **Example**:
  ```go
  type User struct {
      // +govalid:lowercase
      Username string `json:"username"`
  }
  ```
- **Generated Code**:
  ```go
  if !validationhelper.IsLowercase(t.Username) {
      return ErrUsernameLowercaseValidation
  }
  ```

## `govalid:oneof`
- **Description**: Ensures that a field value is one of the specified options (space-separated list). Supports strings and numeric types.
- **Example**:
  ```go
  type Config struct {
      // +govalid:oneof=red green blue
      Color string `json:"color"`
      // +govalid:oneof=1 2 3
      Level int `json:"level"`
  }
  ```
- **Generated Code**:
  ```go
  if !(t.Color == "red" || t.Color == "green" || t.Color == "blue") {
      return ErrColorOneofValidation
  }
  if !(t.Level == 1 || t.Level == 2 || t.Level == 3) {
      return ErrLevelOneofValidation
  }
  ```

## `govalid:number`
- **Description**: Ensures that a string field represents a valid number (including decimals and negative values).
- **Example**:
  ```go
  type Input struct {
      // +govalid:number
      Amount string `json:"amount"`
  }
  ```
- **Generated Code**:
  ```go
  if !validationhelper.IsNumber(t.Amount) {
      return ErrAmountNumberValidation
  }
  ```

## `govalid:alphanum`
- **Description**: Ensures that a string field contains only alphanumeric characters.
- **Example**:
  ```go
  type Product struct {
      // +govalid:alphanum
      SKU string `json:"sku"`
  }
  ```
- **Generated Code**:
  ```go
  if !validationhelper.IsAlphanum(t.SKU) {
      return ErrSKUAlphanumValidation
  }
  ```

## `govalid:containsany`
- **Description**: Ensures that a string field contains at least one of the specified characters.
- **Example**:
  ```go
  type Security struct {
      // +govalid:containsany=!@#$
      Password string `json:"password"`
  }
  ```
- **Generated Code**:
  ```go
  if !strings.ContainsAny(t.Password, "!@#$") {
      return ErrPasswordContainsanyValidation
  }
  ```

## `govalid:excludes`
- **Description**: Ensures that a string field does not contain a specified substring.
- **Example**:
  ```go
  type User struct {
      // +govalid:excludes=admin
      Username string `json:"username"`
  }
  ```
- **Generated Code**:
  ```go
  if strings.Contains(t.Username, "admin") {
      return ErrUsernameExcludesValidation
  }
  ```

## `govalid:excludesall`
- **Description**: Ensures that a string field does not contain any of the specified characters.
- **Example**:
  ```go
  type Comment struct {
      // +govalid:excludesall=<>
      Text string `json:"text"`
  }
  ```
- **Generated Code**:
  ```go
  if strings.ContainsAny(t.Text, "<>") {
      return ErrTextExcludesallValidation
  }
  ```

## `govalid:unique`
- **Description**: Ensures that all elements in a slice are unique.
- **Example**:
  ```go
  type Data struct {
      // +govalid:unique
      Tags []string `json:"tags"`
      // +govalid:unique
      IDs []int `json:"ids"`
  }
  ```
- **Generated Code**:
  ```go
  if func() bool {
      seen := make(map[interface{}]struct{})
      for _, v := range t.Tags {
          if _, exists := seen[v]; exists {
              return true
          }
          seen[v] = struct{}{}
      }
      return false
  }() {
      return ErrTagsUniqueValidation
  }
  ```

## `govalid:uri`
- **Description**: Ensures that a string field is a valid URI (supports various schemes: http, https, ftp, file, etc.).
- **Example**:
  ```go
  type Resource struct {
      // +govalid:uri
      Location string `json:"location"`
  }
  ```
- **Generated Code**:
  ```go
  if !validationhelper.IsValidURI(t.Location) {
      return ErrLocationUriValidation
  }
  ```

## `govalid:fqdn`
- **Description**: Ensures that a string field is a Fully Qualified Domain Name.
- **Example**:
  ```go
  type Server struct {
      // +govalid:fqdn
      Hostname string `json:"hostname"`
  }
  ```
- **Generated Code**:
  ```go
  if !validationhelper.IsValidFQDN(t.Hostname) {
      return ErrHostnameFqdnValidation
  }
  ```

## `govalid:latitude`
- **Description**: Ensures that a string field represents a valid latitude (-90 to 90).
- **Example**:
  ```go
  type Location struct {
      // +govalid:latitude
      Lat string `json:"lat"`
  }
  ```
- **Generated Code**:
  ```go
  if !validationhelper.IsValidLatitude(t.Lat) {
      return ErrLatLatitudeValidation
  }
  ```

## `govalid:longitude`
- **Description**: Ensures that a string field represents a valid longitude (-180 to 180).
- **Example**:
  ```go
  type Location struct {
      // +govalid:longitude
      Lon string `json:"lon"`
  }
  ```
- **Generated Code**:
  ```go
  if !validationhelper.IsValidLongitude(t.Lon) {
      return ErrLonLongitudeValidation
  }
  ```

## `govalid:iscolour` / `govalid:iscolor`
- **Description**: Ensures that a string field represents a valid color. Supports hex (#RGB, #RRGGBB, #RRGGBBAA), rgb/rgba, hsl/hsla, and common named colors.
- **Example**:
  ```go
  type Theme struct {
      // +govalid:iscolour
      Primary string `json:"primary"`
  }
  ```
- **Generated Code**:
  ```go
  if !validationhelper.IsValidColour(t.Primary) {
      return ErrPrimaryIscolourValidation
  }
  ```

## `govalid:minduration`
- **Description**: Ensures that a time.Duration field is at least the specified minimum duration.
- **Example**:
  ```go
  type Config struct {
      // +govalid:minduration=1h
      Timeout time.Duration `json:"timeout"`
  }
  ```
- **Generated Code**:
  ```go
  if func() bool { d, _ := time.ParseDuration("1h"); return t.Timeout < d }() {
      return ErrTimeoutMindurationValidation
  }
  ```

## `govalid:maxduration`
- **Description**: Ensures that a time.Duration field does not exceed the specified maximum duration.
- **Example**:
  ```go
  type Config struct {
      // +govalid:maxduration=24h
      Interval time.Duration `json:"interval"`
  }
  ```
- **Generated Code**:
  ```go
  if func() bool { d, _ := time.ParseDuration("24h"); return t.Interval > d }() {
      return ErrIntervalMaxdurationValidation
  }
  ```

## Conditional Validators

### `govalid:required_if`
- **Description**: Field is required if another field equals a specific value.
- **Format**: `required_if=FieldName Value`
- **Example**:
  ```go
  type Form struct {
      Status string
      // +govalid:required_if=Status active
      ActiveField string `json:"active_field"`
  }
  ```
- **Generated Code**:
  ```go
  if t.Status == "active" && t.ActiveField == "" {
      return ErrActiveFieldRequiredIfValidation
  }
  ```

### `govalid:required_unless`
- **Description**: Field is required unless another field equals a specific value.
- **Format**: `required_unless=FieldName Value`
- **Example**:
  ```go
  type Form struct {
      Status string
      // +govalid:required_unless=Status inactive
      ActiveField string `json:"active_field"`
  }
  ```

### `govalid:required_with`
- **Description**: Field is required when any of the specified fields are present (non-zero).
- **Format**: `required_with=Field1 Field2 ...`
- **Example**:
  ```go
  type Form struct {
      Email string
      // +govalid:required_with=Email
      EmailConfirmation string `json:"email_confirmation"`
  }
  ```

### `govalid:required_with_all`
- **Description**: Field is required when all of the specified fields are present (non-zero).
- **Format**: `required_with_all=Field1 Field2 ...`
- **Example**:
  ```go
  type Form struct {
      FirstName string
      LastName  string
      // +govalid:required_with_all=FirstName LastName
      FullName string `json:"full_name"`
  }
  ```

### `govalid:required_without`
- **Description**: Field is required when any of the specified fields are absent (zero value).
- **Format**: `required_without=Field1 Field2 ...`
- **Example**:
  ```go
  type Contact struct {
      Phone string
      // +govalid:required_without=Phone
      Email string `json:"email"`
  }
  ```

### `govalid:required_without_all`
- **Description**: Field is required when all of the specified fields are absent (zero value).
- **Format**: `required_without_all=Field1 Field2 ...`
- **Example**:
  ```go
  type Contact struct {
      Phone string
      Fax   string
      // +govalid:required_without_all=Phone Fax
      Email string `json:"email"`
  }
  ```

### `govalid:excluded_if`
- **Description**: Field must be absent (zero value) if another field equals a specific value.
- **Format**: `excluded_if=FieldName Value`
- **Example**:
  ```go
  type Form struct {
      Status string
      // +govalid:excluded_if=Status inactive
      InactiveField string `json:"inactive_field"`
  }
  ```

### `govalid:excluded_unless`
- **Description**: Field must be absent unless another field equals a specific value.
- **Format**: `excluded_unless=FieldName Value`

### `govalid:excluded_with`
- **Description**: Field must be absent when any of the specified fields are present.
- **Format**: `excluded_with=Field1 Field2 ...`

### `govalid:excluded_with_all`
- **Description**: Field must be absent when all of the specified fields are present.
- **Format**: `excluded_with_all=Field1 Field2 ...`

### `govalid:excluded_without`
- **Description**: Field must be absent when any of the specified fields are absent.
- **Format**: `excluded_without=Field1 Field2 ...`

### `govalid:excluded_without_all`
- **Description**: Field must be absent when all of the specified fields are absent.
- **Format**: `excluded_without_all=Field1 Field2 ...`

## Summary

govalid now supports **52 validators** covering:
- ✅ Numeric validation (gt, gte, lt, lte, min, eq, ne)
- ✅ String validation (length, pattern, format)
- ✅ Collection validation (size, uniqueness)
- ✅ Format validation (email, URL, UUID, IP, coordinates, colors)
- ✅ Type validation (boolean, numeric, alphanumeric)
- ✅ Duration validation (min/max duration)
- ✅ Conditional validation (12 cross-field validators)
- ✅ Advanced CEL expressions

All validators generate **zero-allocation, type-safe** validation code with comprehensive error messages.
