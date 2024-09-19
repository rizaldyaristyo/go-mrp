package handlers

func permissionIntToString(permissionInt int64) string {
	switch permissionInt {
	case 0:
		return "No Access"
	case 1:
		return "View Only"
	case 2:
		return "Edit Only"
	case 3:
		return "Edit and Approve"
	default:
		return "Invalid Permission"
	}
}
