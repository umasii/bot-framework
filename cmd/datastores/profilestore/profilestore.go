package profilestore

import (
	"fmt"

	store "github.com/umasii/bot-framework/internal/helpers/datastore"
	profiles "github.com/umasii/bot-framework/internal/profiles"

	goErrors "errors"

	errors "github.com/umasii/bot-framework/internal/errors"
)

var profileGroupList []*profiles.ProfileGroup

func AddProfile(profile *profiles.Profile) error {
	profileGroup, err := GetProfileGroupByID(profile.GroupID)
	if err != nil {
		return errors.Handler(goErrors.New("Failed to Add Profile"))
	}
	profile.ProfileID = getNextProfileID(&profileGroup.Profiles)
	profileGroup.Profiles = append(profileGroup.Profiles, *profile)
	store.Write(&profileGroupList, "profiles")
	return nil
}

func AddProfileGroup(profileGroup *profiles.ProfileGroup) {
	profileGroups := GetProfileGroups()
	profileGroup.GroupID = getNextProfileGroupID(&profileGroups)
	if profileGroup.Profiles == nil {
		profileGroup.Profiles = []profiles.Profile{}
	}
	profileGroups = append(profileGroups, *profileGroup)
	store.Write(&profileGroups, "profiles")
}

func GetProfileGroups() []profiles.ProfileGroup {
	var existingProfileGroups []profiles.ProfileGroup
	store.Read(&existingProfileGroups, "profiles")
	return existingProfileGroups
}

func GetProfileByID(groupID, profileID int) (*profiles.Profile, error) {
	profileGroup, err := GetProfileGroupByID(groupID)

	if err != nil {
		return nil, errors.Handler(goErrors.New(fmt.Sprintf("Could not locate Profile Group (GroupID: %v)", groupID)))
	}

	for i := range profileGroup.Profiles {
		if profileGroup.Profiles[i].ProfileID == profileID {
			return &profileGroup.Profiles[i], nil
		}
	}

	return nil, errors.Handler(goErrors.New(fmt.Sprintf("Could not locate Profile (GroupID: %v, ProfileID: %v)", groupID, profileID)))

}

func GetProfileGroupByID(profileGroupID int) (*profiles.ProfileGroup, error) {
	store.Read(&profileGroupList, "profiles")
	for i := range profileGroupList {
		if profileGroupList[i].GroupID == profileGroupID {
			return profileGroupList[i], nil
		}
	}
	return nil, errors.Handler(goErrors.New(fmt.Sprintf("Profile Group with ID %v was not found", profileGroupID)))
}

func getNextProfileID(profiles *[]profiles.Profile) int {
	if len(*profiles) > 0 {
		lastProfile := (*profiles)[len(*profiles)-1]
		return lastProfile.ProfileID + 1
	}
	return 1
}

func getNextProfileGroupID(profileGroups *[]profiles.ProfileGroup) int {
	if len(*profileGroups) > 0 {
		lastProfileGroup := (*profileGroups)[len(*profileGroups)-1]
		return lastProfileGroup.GroupID + 1
	} else {
		return 1
	}
}
