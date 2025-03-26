
DROP TABLE IF EXISTS images;
CREATE TABLE images (
    image_id        INT AUTO_INCREMENT PRIMARY KEY,
    url             VARCHAR(255) NOT NULL UNIQUE,
    name            VARCHAR(255) NULL DEFAULT NULL,
    INDEX idx_images_id (image_id ASC),
    INDEX idx_images_url (url ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

DROP TABLE IF EXISTS siteinfo;
CREATE TABLE siteinfo (
    site_id         BOOL PRIMARY KEY DEFAULT true,
    site_title      VARCHAR(255),
    site_name       VARCHAR(255),
    site_logo_id    INT,
    site_about      TEXT,
    site_copyright  VARCHAR(255),
    contact_address VARCHAR(255),
    contact_email   VARCHAR(255),
    contact_phone   VARCHAR(20),
    CONSTRAINT unique_site_info CHECK (site_id = 1),
    FOREIGN KEY (site_logo_id) REFERENCES images(image_id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

DROP TABLE IF EXISTS users;
CREATE TABLE users (
    user_id         INT AUTO_INCREMENT PRIMARY KEY,
    firstname       VARCHAR(50) NOT NULL,
    lastname        VARCHAR(50) NOT NULL,
    phone           VARCHAR(20) NOT NULL UNIQUE,
    email           VARCHAR(50) NOT NULL UNIQUE,
    password        VARCHAR(255) NOT NULL, -- Assuming hashed
    avatar_id       INT,
    role            ENUM('admin', 'author', 'subscriber') NOT NULL DEFAULT 'subscriber', -- Basic role management
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_users_id (user_id ASC),
    INDEX idx_users_fullname (firstname ASC, lastname ASC),
    INDEX idx_users_phone (phone ASC),
    INDEX idx_users_email (email ASC),
    FULLTEXT(firstname, lastname, phone, email),
    FOREIGN KEY (avatar_id) REFERENCES images(image_id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

DROP TABLE IF EXISTS posts;
CREATE TABLE posts (
    post_id         INT AUTO_INCREMENT PRIMARY KEY,
    user_id         INT, -- Foreign key to the user who authored the post
    title           VARCHAR(255) NOT NULL,
    slug            VARCHAR(255) NOT NULL UNIQUE, -- SEO-friendly URL identifier
    thumbnail_id    INT,
    body            TEXT NOT NULL,
    status          ENUM('draft', 'published') NOT NULL DEFAULT 'draft', -- For content moderation
    private         BOOLEAN NOT NULL DEFAULT 0,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FULLTEXT(title, body),
    UNIQUE INDEX uq_posts_slug (slug ASC),
    INDEX idx_posts_user_id (user_id ASC),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE SET NULL ON UPDATE NO ACTION,
    FOREIGN KEY (thumbnail_id) REFERENCES images(image_id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

DROP TABLE IF EXISTS categories;
CREATE TABLE categories (
    category_id         INT AUTO_INCREMENT PRIMARY KEY,
    parent_category_id  INT NULL DEFAULT NULL,
    name                VARCHAR(100) NOT NULL UNIQUE,
    slug                VARCHAR(255) NOT NULL UNIQUE,

    INDEX idx_parent_category_id (parent_category_id ASC),
    FOREIGN KEY (parent_category_id) REFERENCES categories(category_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

DROP TABLE IF EXISTS post_categories;
CREATE TABLE post_categories (
    post_id         INT NOT NULL,
    category_id     INT NOT NULL,

    INDEX idx_pc_post_id (post_id ASC),
    INDEX idx_pc_category_id (category_id ASC),
    PRIMARY KEY(post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

DROP TABLE IF EXISTS tags;
CREATE TABLE tags (
    tag_id          INT AUTO_INCREMENT PRIMARY KEY,
    name            VARCHAR(50) NOT NULL UNIQUE,
    slug            VARCHAR(100) NOT NULL UNIQUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

DROP TABLE IF EXISTS post_tags;
CREATE TABLE post_tags (
    post_id         INT NOT NULL,
    tag_id          INT NOT NULL,

    INDEX idx_pt_post_id (post_id ASC),
    INDEX idx_pt_tag_id (tag_id ASC),
    PRIMARY KEY(post_id, tag_id),
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(tag_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

-- Chèn dữ liệu mẫu cho bảng images
INSERT INTO images (url, name) VALUES
('https://example.com/logo.png', 'Site Logo'),
('https://example.com/avatar1.jpg', 'Admin Avatar'),
('https://example.com/post1-thumb.jpg', 'Post 1 Thumbnail'),
('https://example.com/post2-thumb.jpg', 'Post 2 Thumbnail');

-- Chèn dữ liệu mẫu cho bảng siteinfo (chỉ 1 bản ghi)
INSERT INTO siteinfo (site_title, site_name, site_logo_id, site_about, site_copyright, contact_address, contact_email, contact_phone)
VALUES
('My Blog', 'Tech Insights', 1, 'A blog about technology and programming', '© 2023 My Blog', '123 Main St, City', 'info@example.com', '123-456-7890');

-- Chèn dữ liệu mẫu cho bảng users
INSERT INTO users (firstname, lastname, phone, email, password, role) VALUES
('TT', 'ASC', '0000', 'ad@ad.ad', '$2a$10$EqoYIVoqP6FeOYbaa2GD7.OtEdyBGUCsue/gvmi5gjxlqL8yi2cg.', 'admin'),
('Jane', 'Smith', '0123456789', 'author@example.com', 'hashed_password_456', NULL, 'author'),
('Bob', 'Johnson', '0369852147', 'subscriber@example.com', 'hashed_password_789', NULL, 'subscriber');

-- Chèn dữ liệu mẫu cho bảng categories (có phân cấp)
INSERT INTO categories (parent_category_id, name, slug) VALUES
(NULL, 'Technology', 'technology'),
(1, 'Programming', 'programming'),
(1, 'Gadgets', 'gadgets'),
(NULL, 'News', 'news'),
(4, 'Tang hoc phi', 'tang-hoc-phi'),
(NULL, 'Announcements', 'announcements'),
(6, 'Thong bao nghi hoc', 'thong-bao-nghi-hoc');

-- Chèn dữ liệu mẫu cho bảng tags
INSERT INTO tags (name, slug) VALUES
('Web Development', 'web-development'),
('Tips & Tricks', 'tips-tricks'),
('Productivity', 'productivity');

-- Chèn dữ liệu mẫu cho bảng posts
INSERT INTO posts (user_id, title, slug, thumbnail_id, body, status, private) VALUES
(2, 'Getting Started with Python', 'getting-started-python', 3, 'Python is a versatile language...', 'published', 0),
(2, 'Top 10 Programming Tools', 'top-10-programming-tools', 4, 'Here are the essential tools...', 'draft', 0),
(3, 'Daily Productivity Hacks', 'daily-productivity-hacks', NULL, 'Boost your efficiency with these tips...', 'published', 0),
(3, 'Private Daily Productivity Hacks', 'private-daily-productivity-hacks', NULL, 'Private Boost your efficiency with these tips...', 'published', 1);

-- Chèn dữ liệu mẫu cho bảng post_categories
INSERT INTO post_categories (post_id, category_id) VALUES
(1, 4), -- Python post → Programming
(2, 4), -- Tools post → Programming
(3, 6), -- Productivity post → Lifestyle
(4, 6); -- Productivity post → Lifestyle

-- Chèn dữ liệu mẫu cho bảng post_tags
INSERT INTO post_tags (post_id, tag_id) VALUES
(1, 1), -- Python post → Web Development
(2, 1), -- Tools post → Web Development
(3, 2), -- Productivity post → Tips & Tricks
(3, 3); -- Productivity post → Productivity
