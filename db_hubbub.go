package readeef

const (
	get_hubbub_subscription    = `SELECT feed_link, lease_duration, verification_time FROM hubbub_subscriptions WHERE link = ?`
	create_hubbub_subscription = `
INSERT INTO hubbub_subscriptions(link, feed_link, lease_duration, verification_time)
	SELECT ?, ?, ?, ? EXCEPT
	SELECT link, feed_link, lease_duration, verification_time
		FROM hubbub_subscriptions WHERE link = ?
`
	update_hubbub_subscription = `
UPDATE hubbub_subscriptions SET feed_link = ?,
	lease_duration = ?, verification_time = ? WHERE link = ?
`
	delete_hubbub_subscription = `DELETE from hubbub_subscriptions where link = ?`
)

func (db DB) GetHubbubSubscription(link string) (HubbubSubscription, error) {
	var s HubbubSubscription

	if err := db.Get(&s, get_hubbub_subscription); err != nil {
		return s, err
	}

	s.Link = link

	return s, nil
}

func (db DB) UpdateHubbubSubscription(s HubbubSubscription) error {
	if err := s.Validate(); err != nil {
		return err
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	ustmt, err := tx.Preparex(update_hubbub_subscription)

	if err != nil {
		return err
	}
	defer ustmt.Close()

	_, err = ustmt.Exec(s.FeedLink, s.LeaseDuration, s.VerificationTime, s.Link)
	if err != nil {
		return err
	}

	cstmt, err := tx.Preparex(create_hubbub_subscription)

	if err != nil {
		return err
	}
	defer cstmt.Close()

	_, err = cstmt.Exec(s.Link, s.FeedLink, s.LeaseDuration, s.VerificationTime, s.Link)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (db DB) DeleteHubbubSubscription(s HubbubSubscription) error {
	if err := s.Validate(); err != nil {
		return err
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex(delete_hubbub_subscription)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(s.Link)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}