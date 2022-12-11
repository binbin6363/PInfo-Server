package service

import (
	"context"
	"log"

	"PInfo-server/api"
)

func (s *Service) GetContactList(ctx context.Context, uid int64) (error, []*api.ContactInfo) {
	err, userContacts := s.dao.GetContactList(ctx, uid)
	if err != nil {
		log.Printf("query contact info failed, err:%v\n", err)
		return err, nil
	}

	var contactList []*api.ContactInfo = nil
	for _, contact := range userContacts {
		contactInfo := &api.ContactInfo{
			Id:           contact.Uid,
			Nickname:     contact.Nickname,
			Gender:       contact.Gender,
			Motto:        contact.Motto,
			Avatar:       contact.Avatar,
			FriendRemark: contact.FriendRemark,
			IsOnline:     0,
		}
		contactList = append(contactList, contactInfo)
	}

	return nil, contactList
}
