-- name: InsertEmail :exec
INSERT INTO Email (subject, expiresAt) 
VALUES (?, ?);

-- name: InsertInbox :exec
INSERT INTO Inbox (id, emailId, address, textContent, htmlContent) 
VALUES (?, ?, ?, ?, ?);

-- name: InsertEmailAddress :exec
INSERT INTO EmailAddress (emailId, type, address) 
VALUES (?, ?, ?);

-- name: GetEmailsForAddress :many
SELECT 
  Inbox.id,
  Email.subject,
  Email.createdAt,
  Email.expiresAt,
  (SELECT address FROM EmailAddress WHERE emailId = Email.id AND type = 'from') as fromAddress,
  (SELECT GROUP_CONCAT(address, ', ') FROM EmailAddress WHERE emailId = Email.id AND type = 'to') as toAddress
FROM Email
JOIN Inbox ON Email.id = Inbox.emailId
WHERE Inbox.address = ?
ORDER BY Email.createdAt DESC;

-- name: GetInboxByID :one
SELECT 
  Inbox.id,
  Inbox.textContent, 
  Inbox.htmlContent, 
  Email.subject, 
  Email.expiresAt,
  Email.createdAt,
  (SELECT address FROM EmailAddress WHERE emailId = Email.id AND type = 'from') as fromAddress,
  (SELECT GROUP_CONCAT(address, ', ') FROM EmailAddress WHERE emailId = Email.id AND type = 'to') as toAddress
FROM Inbox
JOIN Email ON Inbox.emailId = Email.id
WHERE Inbox.id = ?;

-- name: DeleteByInboxID :exec
DELETE FROM Inbox WHERE id = ?;

-- name: DeleteExpiredEntries :exec
DELETE FROM Email WHERE expiresAt < strftime('%s', 'now') * 1000;
