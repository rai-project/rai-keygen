package cmd

func isValidRole(r string) bool {
  r = strings.ToLower(r)
  for _, e := range models.Roles {
    if e == r {
      return true
    }
  }
  return false
}