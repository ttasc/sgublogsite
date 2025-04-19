SET NAMES 'utf8mb4';

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
    contact_phone   VARCHAR(100),
    CONSTRAINT unique_site_info CHECK (site_id = 1),
    FOREIGN KEY (site_logo_id) REFERENCES images(image_id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

DROP TABLE IF EXISTS users;
CREATE TABLE users (
    user_id         INT AUTO_INCREMENT PRIMARY KEY,
    status          ENUM('active', 'inactive') NOT NULL DEFAULT 'active',
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
    FOREIGN KEY (parent_category_id) REFERENCES categories(category_id) ON DELETE SET NULL
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

-- Chèn dữ liệu mẫu cho bảng siteinfo (chỉ 1 bản ghi)
INSERT INTO siteinfo (site_title, site_name, site_logo_id, site_about, site_copyright, contact_address, contact_email, contact_phone)
VALUES
('SGU', 'SGU site', NULL, 'Welcome to SGU site', '© 2025 SGU', ' 273 An Dương Vương – Phường 3 – Quận 5', 'daihocsaigon@sgu.edu.vn', '(84-28) 38.354409 - 38.352309');

-- Chèn dữ liệu mẫu cho bảng users
INSERT INTO users (firstname, lastname, phone, email, password, role) VALUES
('TT', 'ASC', '0000', 'ad@ad.ad', '$2a$10$EqoYIVoqP6FeOYbaa2GD7.OtEdyBGUCsue/gvmi5gjxlqL8yi2cg.', 'admin');

-- Chèn dữ liệu mẫu cho bảng categories (có phân cấp)
INSERT INTO categories (parent_category_id, name, slug) VALUES
(NULL, 'Đào tạo', 'dao-tao'),
(1, 'Giáo dục', 'giao-duc'),
(1, 'Thể thao', 'gadgets'),
(NULL, 'Tin tức', 'news'),
(4, 'Sự kiện', 'su-kien'),
(NULL, 'Thông báo', 'announcements'),
(6, 'Đóng học phí', 'dong-hoc-phi');

-- Chèn dữ liệu mẫu cho bảng tags
INSERT INTO tags (name, slug) VALUES
('học phí', 'hoc-phi'),
('CNTT', 'cntt'),
('sư phạm', 'su-pham');

-- Chèn dữ liệu mẫu cho bảng posts
INSERT INTO posts (user_id, title, slug, thumbnail_id, body, status, private) VALUES
(1, 'Tin tức mới nhất', 'tin-tuc-moi-nhat', NULL, 'Tin tức mới nhất...', 'published', 0),
(1, 'Hội thao toàn thành phố', 'hoi-thao-toan-thanh-pho', NULL, '<p>Nằm trong chuỗi các hoạt động chào mừng các ngày lễ lớn của đất nước
như: kỷ niệm 49 năm Ngày giải phóng miền Nam, thống nhất đất nước
(30/4/1975 – 30/4/2024); 138 năm Ngày Quốc tế Lao động (01/5/1886 –
01/5/2024); 134 năm Ngày sinh Chủ tịch Hồ Chí Minh (19/5/1890 –
19/5/2024); hưởng ứng Tháng Công nhân lần thứ 16 và Tháng hành động về
an toàn, vệ sinh lao động năm 2024; Hội thao truyền thống Công đoàn
Trường Đại học Sài Gòn năm học 2023-2024 đã được tổ chức từ ngày
22/4/2024 đến ngày 11/5/2024.</p>



<p>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; Hội thao gồm 4 môn thi đấu: Bóng bàn, Bóng chuyền, Cầu lông,
Cờ tướng, đã thu hút sự tham gia của hơn 270 lượt vận động viên đến từ
36 Công đoàn bộ phận.</p>



<p>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; Là hoạt động thường niên của Công đoàn, Hội thao truyền thống
 đã tạo sân chơi lành mạnh, bổ ích; tạo cơ hội giao lưu, luyện tập thể
dục thể thao; góp phần rèn luyện sức khoẻ, tăng cường tình đoàn kết giữa
 các viên chức, người lao động của Trường.</p>



<p>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; Các vận động viên (VĐV) đã thi đấu với tinh thần quyết liệt, vui tươi, phấn khởi và đạt nhiều thành tích xuất sắc.</p>



<p>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <strong>1. Môn Bóng chuyền:</strong></p>



<p>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <strong><em>Giải Nhất:</em></strong> Đội B (Liên quân các
CĐBP: Phòng Đào tạo Sau Đại học, Phòng Khảo thí và Đảm bảo chất lượng
giáo dục, Khoa Nghệ thuật, Khoa Ngoại ngữ, Khoa Thư viện Văn phòng, Khoa
 Điện tử viễn thông, Khoa Giáo dục Quốc phòng An ninh và Giáo dục thể
chất)</p>



<p>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <strong><em>Giải Nhì:</em></strong> Đội A (Liên quân các
CĐBP: Văn phòng, Phòng Giáo dục thường xuyên, Phòng Công tác sinh viên,
Phòng Kế hoạch Tài chính, Trung tâm Công nghệ thông tin, Khoa Toán – Ứng
 dụng).</p>



<p>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <strong><em>Giải Ba: </em></strong>Đội D (Liên quân CĐBP Ban Quản lý Dự án và Hạ tầng và Trung tâm Ngoại ngữ).</p>



<p>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <strong>2.</strong> <strong>Môn Cờ tướng:</strong></p>



<figure class="wp-block-table">
<table>
<tbody>
<tr>
<td><strong>Giải</strong></td>
<td><strong>Họ tên VĐV</strong></td>
<td><strong>Công đoàn bộ phận</strong></td>
</tr>
<tr>
<td>Nhất</td>
<td>Nguyễn Hữu Duy Khang</td>
<td>Khoa Sư phạm Khoa học Tự nhiên</td>
</tr>
<tr>
<td>Nhì</td>
<td>Phan Đức Tuấn</td>
<td>Khoa Toán – Ứng dụng</td>
</tr>
<tr>
<td>Ba</td>
<td>Trần Thanh Hiệp</td>
<td>Ban Quản lý Dự án và Hạ tầng</td>
</tr>
<tr>
<td>Ba</td>
<td>Nguyễn Văn Mười</td>
<td>Ban Quản lý Dự án và Hạ tầng</td>
</tr>
</tbody>
</table>
</figure>



<p>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <strong>3.</strong> <strong>Môn Bóng bàn:</strong></p>



<figure class="wp-block-table">
<table>
<tbody>
<tr>
<td><strong>Nội dung</strong></td>
<td><strong>Giải </strong><strong>Nhất</strong></td>
<td><strong>Giải </strong><strong>Nhì</strong></td>
<td><strong>Giải </strong><strong>Ba</strong></td>
</tr>
<tr>
<td><strong><em>Đơn Nam</em></strong></td>
<td>Nguyễn Hoàng Tuấn (TTHL)</td>
<td>Bùi Đức Tú (K.GD)</td>
<td>1. Trần Thắng Thành (P.CTSV)<br/>2. Huỳnh Lê Minh Thiện (K.ĐTVT)</td>
</tr>
<tr>
<td><strong><em>Đơn Nữ</em></strong></td>
<td>Nguyễn Phan Thu Hằng (K.QTKD)</td>
<td>Lê Thị Khánh Vân (P.ĐT</td>
<td>1. Nguyễn Trịnh Tố Anh (K.NN)<br/>2. Đoàn Ngọc Anh (P.ĐT)</td>
</tr>
<tr>
<td><strong><em>Đôi Nam</em></strong></td>
<td>Lý Hoàng Ánh (K.QTKD) <br/>Lê Mai Hải (K.QTKD)</td>
<td>Trần Thắng Thành (P.CTSV) <br/>Huỳnh Lê Minh Thiện (K.ĐTVT)</td>
<td>Lê Tùng Lâm (K.VH&amp;DL) <br/>Trương Lập Trọng Thủy (TYT)</td>
</tr>
<tr>
<td><strong><em>Đôi</em></strong> <strong><em>Nam – Nữ</em></strong></td>
<td>Lý Hoàng Ánh (K.QTKD) <br/>Nguyễn Phan Thu Hằng (K.QTKD)</td>
<td>Nguyễn Hoàng Tuấn (TTHL) <br/>Đỗ Nguyễn Thanh Trúc (THSG)</td>
<td>1. Phạm Hoàng Vương (K.CNTT) <br/>Đoàn Ngọc Anh (P.ĐT) <br/>2. Nguyễn An Hòa (K.GD) <br/>Nguyễn Trịnh Tố Anh (K.NN)</td>
</tr>
<tr>
<td><strong><em>Đơn Nam nâng cao</em></strong></td>
<td>Phạm Hoàng Vương (K.CNTT)</td>
<td>Hồ Cảnh Hoàng Giang (K.SPKHTN)</td>
<td>1. Mỵ Giang Sơn (K.GD) <br/>2. Nguyễn Đức Hưng (K.SPKHTN)</td>
</tr>
<tr>
<td><strong><em>Đôi Nam</em></strong> <strong><em>nâng cao</em></strong></td>
<td>Hồ Cảnh Hoàng Giang (K.SPKHTN) <br/>Nguyễn Đức Hưng (K.SPKHTN)</td>
<td>Phạm Hoàng Vương (K.CNTT) <br/>Đặng Xuân Dự (K.SPKHTN)</td>
<td>Mỵ Giang Sơn (K.GD) <br/>Bùi Đức Tú (K.GD)</td>
</tr>
</tbody>
</table>
</figure>



<p>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <strong>4.</strong> <strong>Môn Cầu lông:</strong></p>



<figure class="wp-block-table">
<table>
<thead>
<tr>
<td><strong>Nội dung</strong></td>
<td><strong>Giải </strong><strong>Nhất</strong></td>
<td><strong>Giải </strong><strong>Nhì</strong></td>
<td><strong>Giải </strong><strong>Ba</strong></td>
</tr>
</thead>
<tbody>
<tr>
<td><strong><em>Đơn Nam</em></strong></td>
<td>Huỳnh Thanh Hiếu (P.QLKH)</td>
<td>Đoàn Thế Vinh (B.QLDA&amp;HT)</td>
<td>1. Lê Quốc Dũng (THSG) <br/>2. Vũ Hùng Phi (P.TTPC)</td>
</tr>
<tr>
<td><strong><em>Đôi Nam</em></strong></td>
<td>Đỗ Cảnh Phụng + Lê Quốc Dũng (THSG)</td>
<td>Giang Quốc Tuấn + Lê Nhân Tâm (P.KHTC)</td>
<td>1. Hồ Đình Khuê + Vũ Hùng Phi (P.TTPC) <br/>2. Đoàn Thế Vinh + Huỳnh Thanh Trung (B.QLDA&amp;HT)</td>
</tr>
<tr>
<td><strong><em>Đôi Nữ</em></strong></td>
<td>Nguyễn Thu Thủy + Nguyễn Thị Nở (P.TCCB)</td>
<td>Trịnh Thị Huyền Thương + Nguyễn Ngọc Huyền Trân (K.TCKT)</td>
<td>1. Nguyễn Thị Hương + Lê Thị Việt Kiều (TT.TC&amp;SK) <br/>2. Nguyễn Thị Thanh Trúc + Phạm Thị Nga (THSG)</td>
</tr>
<tr>
<td><strong><em>Đôi</em></strong> <strong><em>Nam – Nữ</em></strong></td>
<td>Nguyễn Hữu Trọng + Võ Bạch Minh Thy (TT.CNTT)</td>
<td>Đỗ Cảnh Phụng + Đỗ Thị Thanh Trúc (THSG)</td>
<td>1. Giang Quốc Tuấn + Nguyễn Kim Dung (P.KHTC) <br/>2. Phan Anh Huy + Nguyễn Thu Thủy (P.TCCB)</td>
</tr>
<tr>
<td><strong><em>Đôi N</em></strong><strong><em>am nâng cao</em></strong></td>
<td>Đặng Trung Nam (K.NN) <br/>Lê Hoàng Dũng (TTHL)</td>
<td>Đỗ Quang Tuấn (TT.CNTT) <br/>Huỳnh Vạng Phước (K.GDQP-AN&amp;GDTC)</td>
<td>Nguyễn Hữu Trí (K.SPKHTN) <br/>Nguyễn Thành Trung (THSG)</td>
</tr>
<tr>
<td><strong><em>Đôi Nữ nâng cao</em></strong></td>
<td>Nguyễn Thị Thanh Thảo (P.KHTC) <br/>Huỳnh Xuân Trúc (Tiểu học THĐHSG)</td>
<td>Nguyễn Ngọc Dung (K.MT) <br/>Vương Thảo Nguyên (K.NN)</td>
<td>Dương Thị Thu Vân (P.ĐT) <br/>Hoàng Thị Phương Thúy (K.TVVP)</td>
</tr>
<tr>
<td><strong><em>Đôi</em></strong> <strong><em>Nam – Nữ</em></strong> <strong><em>nâng cao</em></strong></td>
<td>Nguyễn Thành Trung (THSG) <br/>Nguyễn Thị Thanh Thảo (P.KHTC)</td>
<td>Đỗ Quang Tuấn (TT.CNTT) <br/>Lê Thu Hiền (K.GDQP-AN&amp;GDTC)</td>
<td>1. Đặng Trung Nam (K.NN) <br/>Huỳnh Xuân Trúc (Tiểu học THĐHSG) <br/>2. Lê Hoàng Dũng (TTHL) <br/>Dương Thị Thu Vân (P.ĐT)</td>
</tr>
</tbody>
</table>
</figure>



<p>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <strong>Một số hình ảnh</strong></p>



<figure class="wp-block-image size-full"><img class="alignnone wp-image-2612 size-full lazy-load-active" src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_5-scaled.jpg" data-src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_5-scaled.jpg" alt="" width="2560" height="1537" data-srcset="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_5-scaled.jpg 2560w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_5-350x210.jpg 350w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_5-768x461.jpg 768w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_5-1536x922.jpg 1536w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_5-2048x1229.jpg 2048w"/></figure>



<figure class="wp-block-image aligncenter size-full">
<figure id="attachment_2613" aria-describedby="caption-attachment-2613" style="width: 2560px" class="wp-caption alignnone"><img class="wp-image-2613 lazy-load-active" src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_1.jpg" data-src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_1.jpg" alt="" width="2560" height="1920" data-srcset="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_1.jpg 2560w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_1-350x263.jpg 350w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_1-768x576.jpg 768w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_1-1536x1152.jpg 1536w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_1-2048x1536.jpg 2048w"/><figcaption id="caption-attachment-2613" class="wp-caption-text">VĐV thi đấu Bóng bàn</figcaption></figure>
<figcaption class="wp-element-caption"></figcaption>
</figure>





<figure class="wp-block-image aligncenter size-full">
<figcaption class="wp-element-caption">
<figure id="attachment_2788" aria-describedby="caption-attachment-2788" style="width: 2560px" class="wp-caption alignnone"><img class="wp-image-2788 size-full lazy-load-active" src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT14-scaled.jpg" data-src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT14-scaled.jpg" alt="" width="2560" height="1707" data-srcset="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT14-scaled.jpg 2560w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT14-350x233.jpg 350w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT14-768x512.jpg 768w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT14-1536x1024.jpg 1536w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT14-2048x1365.jpg 2048w"/><figcaption id="caption-attachment-2788" class="wp-caption-text">VĐV thi đấu Bóng chuyền</figcaption></figure>
<figure id="attachment_2789" aria-describedby="caption-attachment-2789" style="width: 2560px" class="wp-caption alignnone"><img class="wp-image-2789 size-full lazy-load-active" src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_3-scaled.jpg" data-src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_3-scaled.jpg" alt="" width="2560" height="1920" data-srcset="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_3-scaled.jpg 2560w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_3-350x263.jpg 350w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_3-768x576.jpg 768w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_3-1536x1152.jpg 1536w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_3-2048x1536.jpg 2048w"/><figcaption id="caption-attachment-2789" class="wp-caption-text">VĐV thi đấu Cầu lông</figcaption></figure>
</figcaption>
<figure id="attachment_2790" aria-describedby="caption-attachment-2790" style="width: 2560px" class="wp-caption alignnone"><img class="wp-image-2790 size-full lazy-load-active" src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_4-scaled.jpg" data-src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_4-scaled.jpg" alt="" width="2560" height="1707" data-srcset="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_4-scaled.jpg 2560w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_4-350x233.jpg 350w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_4-768x512.jpg 768w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_4-1536x1024.jpg 1536w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/DHSG_4-2048x1365.jpg 2048w"/><figcaption id="caption-attachment-2790" class="wp-caption-text">VĐV thi đấu Cờ tướng</figcaption></figure>
</figure>













<figure class="wp-block-image aligncenter size-full">
<figcaption class="wp-element-caption">
<figure id="attachment_2791" aria-describedby="caption-attachment-2791" style="width: 2560px" class="wp-caption aligncenter"><img class="wp-image-2791 size-full lazy-load-active" src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT7-scaled.jpg" data-src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT7-scaled.jpg" alt="" width="2560" height="1706" data-srcset="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT7-scaled.jpg 2560w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT7-350x233.jpg 350w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT7-768x512.jpg 768w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT7-1536x1024.jpg 1536w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT7-2048x1365.jpg 2048w"/><figcaption id="caption-attachment-2791" class="wp-caption-text">PGS.TS. Phạm Hoàng Quân – Hiệu trưởng trao giải cho các VĐV thi đấu môn Bóng bàn</figcaption></figure>
</figcaption>
<figure id="attachment_2792" aria-describedby="caption-attachment-2792" style="width: 2560px" class="wp-caption aligncenter"><img class="wp-image-2792 size-full lazy-load-active" src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT9-scaled.jpg" data-src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT9-scaled.jpg" alt="" width="2560" height="1707" data-srcset="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT9-scaled.jpg 2560w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT9-350x233.jpg 350w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT9-768x512.jpg 768w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT9-1536x1024.jpg 1536w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/HT9-2048x1365.jpg 2048w"/><figcaption id="caption-attachment-2792" class="wp-caption-text">PGS.TS. Lê Chi Lan – Phó Hiệu trưởng trao giải cho các VĐV thi đấu môn Bóng bàn</figcaption></figure>
<figure id="attachment_2793" aria-describedby="caption-attachment-2793" style="width: 2560px" class="wp-caption alignnone"><figcaption id="caption-attachment-2793" class="wp-caption-text">TS. Hồ Kỳ Quang Minh – Bí thư Đảng ủy, Chủ tịch Hội đồng Trường trao giải cho các VĐV thi đấu môn Bóng chuyền</figcaption></figure>
<figure id="attachment_2794" aria-describedby="caption-attachment-2794" style="width: 2560px" class="wp-caption alignnone"><img class="wp-image-2794 size-full lazy-load-active" src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/CT2-scaled.jpg" data-src="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/CT2-scaled.jpg" alt="" width="2560" height="1706" data-srcset="https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/CT2-scaled.jpg 2560w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/CT2-350x233.jpg 350w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/CT2-768x512.jpg 768w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/CT2-1536x1024.jpg 1536w, https://congdoan.sgu.edu.vn/wp-content/uploads/2024/05/CT2-2048x1365.jpg 2048w"/><figcaption id="caption-attachment-2794" class="wp-caption-text">ThS. Lê Chí Cường – Chủ tịch Công đoàn Trường trao giải cho các VĐV thi đấu môn Cờ tướng</figcaption></figure></figure>', 'published', 0),
(1, 'Thông báo đóng học phí', 'thong-bao-dong-hoc-phi', NULL, 'Thống báo đóng học phí...', 'published', 0),
(1, 'Lịch công tác', 'lich-cong-tac', NULL, 'Lịch công tác...', 'published', 1);

-- Chèn dữ liệu mẫu cho bảng post_categories
INSERT INTO post_categories (post_id, category_id) VALUES
(1, 4),
(2, 4),
(3, 6),
(4, 6);

-- Chèn dữ liệu mẫu cho bảng post_tags
INSERT INTO post_tags (post_id, tag_id) VALUES
(1, 1),
(2, 1),
(3, 2),
(3, 3);
