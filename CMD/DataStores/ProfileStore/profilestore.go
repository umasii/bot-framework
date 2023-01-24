package ProfileStore

import (
	"errors"
	"fmt"

	store "github.com/cicadaaio/LVBot/Internal/Helpers/DataStore"
	"github.com/cicadaaio/LVBot/Internal/Profiles"

	"github.com/cicadaaio/LVBot/Internal/Errors"
)

var profileGroupList []*Profiles.ProfileGroup

func AddProfile(profile *Profiles.Profile) error {
	profileGroup, err := GetProfileGroupByID(profile.GroupID)
	if err != nil {
		return Errors.Handler(errors.New("Failed to Add Profile"))
	}
	profile.ProfileID = getNextProfileID(&profileGroup.Profiles)
	profileGroup.Profiles = append(profileGroup.Profiles, *profile)
	store.Write(&profileGroupList, "profiles")
	return nil
}

func AddProfileGroup(profileGroup *Profiles.ProfileGroup) {
	profileGroups := GetProfileGroups()
	profileGroup.GroupID = getNextProfileGroupID(&profileGroups)
	if profileGroup.Profiles == nil {
		profileGroup.Profiles = []Profiles.Profile{}
	}
	profileGroups = append(profileGroups, *profileGroup)
	store.Write(&profileGroups, "profiles")
}

func GetProfileGroups() []Profiles.ProfileGroup {
	var existingProfileGroups []Profiles.ProfileGroup
	store.Read(&existingProfileGroups, "profiles")
	return existingProfileGroups
}

func GetProfileByID(groupID, profileID int) (*Profiles.Profile, error) {
	profileGroup, err := GetProfileGroupByID(groupID)

	if err != nil {
		return nil, Errors.Handler(errors.New(fmt.Sprintf("Could not locate Profile Group (GroupID: %v)", groupID)))
	}

	for i := range profileGroup.Profiles {
		if profileGroup.Profiles[i].ProfileID == profileID {
			return &profileGroup.Profiles[i], nil
		}
	}

	return nil, Errors.Handler(errors.New(fmt.Sprintf("Could not locate Profile (GroupID: %v, ProfileID: %v)", groupID, profileID)))

}

func GetProfileGroupByID(profileGroupID int) (*Profiles.ProfileGroup, error) {
	store.Read(&profileGroupList, "profiles")
	for i := range profileGroupList {
		if profileGroupList[i].GroupID == profileGroupID {
			return profileGroupList[i], nil
		}
	}
	return nil, Errors.Handler(errors.New(fmt.Sprintf("Profile Group with ID %v was not found", profileGroupID)))
}

func getNextProfileID(profiles *[]Profiles.Profile) int {
	if len(*profiles) > 0 {
		lastProfile := (*profiles)[len(*profiles)-1]
		return lastProfile.ProfileID + 1
	}
	return 1
}

func getNextProfileGroupID(profileGroups *[]Profiles.ProfileGroup) int {
	if len(*profileGroups) > 0 {
		lastProfileGroup := (*profileGroups)[len(*profileGroups)-1]
		return lastProfileGroup.GroupID + 1
	} else {
		return 1
	}
}
