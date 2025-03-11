
DROP TABLE IF EXISTS images;
CREATE TABLE images (
    image_id        INT AUTO_INCREMENT PRIMARY KEY,
    url             VARCHAR(255) NOT NULL UNIQUE,
    name            VARCHAR(255) NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

DROP TABLE IF EXISTS users;
CREATE TABLE users (
    user_id         INT AUTO_INCREMENT PRIMARY KEY,
    firstname       VARCHAR(50) NOT NULL,
    lastname        VARCHAR(50) NOT NULL,
    mobile          VARCHAR(20) NOT NULL UNIQUE,
    email           VARCHAR(50) NOT NULL UNIQUE,
    password        VARCHAR(255) NOT NULL, -- Assuming hashed
    profile_pic_id  INT,
    role            ENUM('admin', 'author', 'subscriber') DEFAULT 'subscriber', -- Basic role management
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_users_id (user_id ASC),
    INDEX idx_users_fullname (firstname ASC, lastname ASC),
    INDEX idx_users_mobile (mobile ASC),
    INDEX idx_users_email (email ASC),
    FULLTEXT(firstname, lastname, mobile, email),
    FOREIGN KEY (profile_pic_id) REFERENCES images(image_id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

DROP TABLE IF EXISTS posts;
CREATE TABLE posts (
    post_id         INT AUTO_INCREMENT PRIMARY KEY,
    user_id         INT, -- Foreign key to the user who authored the post
    title           VARCHAR(255) NOT NULL,
    slug            VARCHAR(255) NOT NULL UNIQUE, -- SEO-friendly URL identifier
    preview_pic_id  INT,
    body            TEXT NOT NULL,
    status          ENUM('draft', 'published') NOT NULL DEFAULT 'draft', -- For content moderation
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FULLTEXT(title, body),
    UNIQUE INDEX uq_posts_slug (slug ASC),
    INDEX idx_posts_user_id (user_id ASC),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE SET NULL ON UPDATE NO ACTION,
    FOREIGN KEY (preview_pic_id) REFERENCES images(image_id) ON DELETE SET NULL
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

