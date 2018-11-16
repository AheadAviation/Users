// Copyright © 2018 Tim Curless <tim.curless@thinkahead.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongodb

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/aheadaviation/Users/users"
)

var (
	name            string
	password        string
	host            string
	db              = "users"
	ErrInvalidHexID = errors.New("Invalid Id Hex")
)

func init() {
	flag.StringVar(&name, "mongo-user", os.Getenv("MONGO_USER"), "Mongo Username")
	flag.StringVar(&password, "mongo-password", os.Getenv("MONGO_PASS"), "Mongo Password")
	flag.StringVar(&host, "mongo-host", os.Getenv("MONGO_HOST"), "Mongo Host")
}

type Mongo struct {
	Session *mgo.Session
}

func (m *Mongo) Init() error {
	u := getURL()
	var err error
	m.Session, err = mgo.DialWithTimeout(u.String(), time.Duration(5)*time.Second)
	if err != nil {
		return nil
	}
	return m.EnsureIndexes()
}

type MongoUser struct {
	users.User `bson:",inline"`
	ID         bson.ObjectId   `bson:"_id"`
	AddressIDs []bson.ObjectId `bson:"addresses"`
	CardIDs    []bson.ObjectId `bson:"cards"`
}

func New() MongoUser {
	u := users.New()
	return MongoUser{
		User:       u,
		AddressIDs: make([]bson.ObjectId, 0),
		CardIDs:    make([]bson.ObjectId, 0),
	}
}

func (mu *MongoUser) AddUserIds() {
	if mu.User.Addresses == nil {
		mu.User.Addresses = make([]users.Address, 0)
	}
	for _, id := range mu.AddressIDs {
		mu.User.Addresses = append(mu.User.Addresses, users.Address{
			ID: id.Hex(),
		})
	}
	if mu.User.Cards == nil {
		mu.User.Cards = make([]users.Card, 0)
	}
	for _, id := range mu.CardIDs {
		mu.User.Cards = append(mu.User.Cards, users.Card{ID: id.Hex()})
	}
	mu.User.UserID = mu.ID.Hex()
}

type MongoAddress struct {
	users.Address `bson:",inline"`
	ID            bson.ObjectId `bson:"_id"`
}

func (m *MongoAddress) AddID() {
	m.Address.ID = m.ID.Hex()
}

type MongoCard struct {
	users.Card `bson:",inline"`
	ID         bson.ObjectId `bson:"_id"`
}

func (m *MongoCard) AddID() {
	m.Card.ID = m.ID.Hex()
}

func (m *Mongo) CreateUser(u *users.User) error {
	s := m.Session.Copy()
	defer s.Close()
	id := bson.NewObjectId()
	mu := New()
	mu.User = *u
	mu.ID = id
	var carderr error
	var addrerr error
	mu.CardIDs, carderr = m.createCards(u.Cards)
	mu.AddressIDs, addrerr = m.createAddresses(u.Addresses)
	c := s.DB("").C("customers")
	_, err := c.UpsertId(mu.ID, mu)
	if err != nil {
		m.cleanAttributes(mu)
		return err
	}
	mu.User.UserID = mu.ID.Hex()
	if carderr != nil || addrerr != nil {
		return fmt.Errorf("%v %v", carderr, addrerr)
	}
	*u = mu.User
	return nil
}

func (m *Mongo) createCards(cs []users.Card) ([]bson.ObjectId, error) {
	s := m.Session.Copy()
	defer s.Close()
	ids := make([]bson.ObjectId, 0)
	defer s.Close()
	for k, ca := range cs {
		id := bson.NewObjectId()
		mc := MongoCard{Card: ca, ID: id}
		c := s.DB("").C("cards")
		_, err := c.UpsertId(mc.ID, mc)
		if err != nil {
			return ids, err
		}
		ids = append(ids, id)
		cs[k].ID = id.Hex()
	}
	return ids, nil
}

func (m *Mongo) createAddresses(as []users.Address) ([]bson.ObjectId, error) {
	ids := make([]bson.ObjectId, 0)
	s := m.Session.Copy()
	defer s.Close()
	for k, a := range as {
		id := bson.NewObjectId()
		ma := MongoAddress{Address: a, ID: id}
		c := s.DB("").C("addresses")
		_, err := c.UpsertId(ma.ID, ma)
		if err != nil {
			return ids, err
		}
		ids = append(ids, id)
		as[k].ID = id.Hex()
	}
	return ids, nil
}

func (m *Mongo) cleanAttributes(mu MongoUser) error {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("addresses")
	_, err := c.RemoveAll(bson.M{"_id": bson.M{"$in": mu.AddressIDs}})
	c = s.DB("").C("cards")
	_, err = c.RemoveAll(bson.M{"_id": bson.M{"$in": mu.CardIDs}})
	return err
}

func (m *Mongo) appendAttributeId(attr string, id bson.ObjectId, userid string) error {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("customers")
	return c.Update(bson.M{"_id": bson.ObjectIdHex(userid)},
		bson.M{"$addToSet": bson.M{attr: id}})
}

func (m *Mongo) removeAttributeId(attr, userid string, id bson.ObjectId) error {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("customers")
	return c.Update(bson.M{"_id": bson.ObjectIdHex(userid)},
		bson.M{"$pull": bson.M{attr: id}})
}

func (m *Mongo) GetUserByName(name string) (users.User, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("customers")
	mu := New()
	err := c.Find(bson.M{"username": name}).One(&mu)
	mu.AddUserIds()
	return mu.User, err
}

func (m *Mongo) GetUser(id string) (users.User, error) {
	s := m.Session.Copy()
	defer s.Close()
	if !bson.IsObjectIdHex(id) {
		return users.New(), errors.New("Invalid id hex")
	}
	c := s.DB("").C("customers")
	mu := New()
	err := c.FindId(bson.ObjectIdHex(id)).One(&mu)
	mu.AddUserIds()
	return mu.User, err
}

func (m *Mongo) GetUsers() ([]users.User, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("customers")
	var mus []MongoUser
	err := c.Find(nil).All(&mus)
	us := make([]users.User, 0)
	for _, mu := range mus {
		mu.AddUserIds()
		us = append(us, mu.User)
	}
	return us, err
}

func (m *Mongo) GetUserAttributes(u *users.User) error {
	s := m.Session.Copy()
	defer s.Close()
	ids := make([]bson.ObjectId, 0)
	for _, a := range u.Addresses {
		if !bson.IsObjectIdHex(a.ID) {
			return ErrInvalidHexID
		}
		ids = append(ids, bson.ObjectIdHex(a.ID))
	}
	var ma []MongoAddress
	c := s.DB("").C("addresses")
	err := c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&ma)
	if err != nil {
		return err
	}
	na := make([]users.Address, 0)
	for _, a := range ma {
		a.Address.ID = a.ID.Hex()
		na = append(na, a.Address)
	}
	u.Addresses = na

	ids = make([]bson.ObjectId, 0)
	for _, c := range u.Cards {
		if !bson.IsObjectIdHex(c.ID) {
			return ErrInvalidHexID
		}
		ids = append(ids, bson.ObjectIdHex(c.ID))
	}
	var mc []MongoCard
	c = s.DB("").C("cards")
	err = c.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&mc)
	if err != nil {
		return err
	}

	nc := make([]users.Card, 0)
	for _, ca := range mc {
		ca.Card.ID = ca.ID.Hex()
		nc = append(nc, ca.Card)
	}
	u.Cards = nc
	return nil
}

func (m *Mongo) GetCard(id string) (users.Card, error) {
	s := m.Session.Copy()
	defer s.Close()
	if !bson.IsObjectIdHex(id) {
		return users.Card{}, errors.New("Invalid id hex")
	}
	c := s.DB("").C("cards")
	mc := MongoCard{}
	err := c.FindId(bson.ObjectIdHex(id)).One(&mc)
	mc.AddID()
	return mc.Card, err
}

func (m *Mongo) GetCards() ([]users.Card, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("cards")
	var mcs []MongoCard
	err := c.Find(nil).All(&mcs)
	cs := make([]users.Card, 0)
	for _, mc := range mcs {
		mc.AddID()
		cs = append(cs, mc.Card)
	}
	return cs, err
}

func (m *Mongo) CreateCard(ca *users.Card, userid string) error {
	if userid != "" && !bson.IsObjectIdHex(userid) {
		return errors.New("Invalid id hex")
	}
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("cards")
	id := bson.NewObjectId()
	mc := MongoCard{Card: *ca, ID: id}
	_, err := c.UpsertId(mc.ID, mc)
	if err != nil {
		return err
	}

	if userid != "" {
		err = m.appendAttributeId("cards", mc.ID, userid)
		if err != nil {
			return err
		}
	}
	mc.AddID()
	*ca = mc.Card
	return err
}

func (m *Mongo) GetAddress(id string) (users.Address, error) {
	s := m.Session.Copy()
	defer s.Close()
	if !bson.IsObjectIdHex(id) {
		return users.Address{}, errors.New("Invalid id hex")
	}
	c := s.DB("").C("addresses")
	ma := MongoAddress{}
	err := c.FindId(bson.ObjectIdHex(id)).One(&ma)
	ma.AddID()
	return ma.Address, err
}

func (m *Mongo) GetAddresses() ([]users.Address, error) {
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("addresses")
	var mas []MongoAddress
	err := c.Find(nil).All(&mas)
	as := make([]users.Address, 0)
	for _, ma := range mas {
		ma.AddID()
		as = append(as, ma.Address)
	}
	return as, err
}

func (m *Mongo) CreateAddress(a *users.Address, userid string) error {
	if userid != "" && !bson.IsObjectIdHex(userid) {
		return errors.New("Invalid id hex")
	}
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C("addresses")
	id := bson.NewObjectId()
	ma := MongoAddress{Address: *a, ID: id}
	_, err := c.UpsertId(ma.ID, ma)
	if err != nil {
		return err
	}

	if userid != "" {
		err = m.appendAttributeId("addresses", ma.ID, userid)
		if err != nil {
			return err
		}
	}
	ma.AddID()
	*a = ma.Address
	return err
}

func (m *Mongo) Delete(entity, id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("invalid id hex")
	}
	s := m.Session.Copy()
	defer s.Close()
	c := s.DB("").C(entity)
	if entity == "customers" {
		u, err := m.GetUser(id)
		if err != nil {
			return err
		}
		aids := make([]bson.ObjectId, 0)
		for _, a := range u.Addresses {
			aids = append(aids, bson.ObjectIdHex(a.ID))
		}
		cids := make([]bson.ObjectId, 0)
		for _, c := range u.Cards {
			cids = append(cids, bson.ObjectIdHex(c.ID))
		}
		ac := s.DB("").C("addresses")
		ac.RemoveAll(bson.M{"_id": bson.M{"$in": aids}})
		cc := s.DB("").C("cards")
		cc.RemoveAll(bson.M{"_id": bson.M{"$in": cids}})
	} else {
		c := s.DB("").C("customers")
		c.UpdateAll(bson.M{},
			bson.M{"$pull": bson.M{entity: bson.ObjectIdHex(id)}})
	}
	return c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}

func (m *Mongo) EnsureIndexes() error {
	s := m.Session.Copy()
	defer s.Close()
	i := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     false,
	}
	c := s.DB("").C("customers")
	return c.EnsureIndex(i)
}

func (m *Mongo) Ping() error {
	s := m.Session.Copy()
	defer s.Close()
	return s.Ping()
}

func getURL() url.URL {
	ur := url.URL{
		Scheme: "mongodb",
		Host:   host,
		Path:   db,
	}
	if name != "" {
		u := url.UserPassword(name, password)
		ur.User = u
	}
	return ur
}
