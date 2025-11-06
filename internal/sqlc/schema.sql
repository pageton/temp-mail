CREATE TABLE IF NOT EXISTS Email (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  subject TEXT,
  createdAt INTEGER DEFAULT (strftime('%s', 'now') * 1000),
  expiresAt DATETIME 
);

CREATE TABLE IF NOT EXISTS Inbox (
  id TEXT PRIMARY KEY,
  address TEXT,
  textContent TEXT,
  htmlContent TEXT,
  createdAt INTEGER DEFAULT (strftime('%s', 'now') * 1000),

  emailId INTEGER,
  FOREIGN KEY (emailId) REFERENCES Email(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS EmailAddress (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  type TEXT, -- Can be 'from', 'to'
  address TEXT,

  emailId INTEGER,
  FOREIGN KEY (emailId) REFERENCES Email(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_email_id ON EmailAddress(emailId);
CREATE INDEX IF NOT EXISTS idx_inbox_address ON Inbox(address);
