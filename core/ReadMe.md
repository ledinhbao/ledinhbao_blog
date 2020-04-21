# Package Core

## Mô tả

Core bao gồm tất cả những thành phần cơ bản mà tất cả các module khác phụ thuộc vào:

- `User`: lưu thông tin username, password, email, rank.
- `Connection`: tạo liên kết tới database tùy theo dialect được chỉ định.
- `Config`: đọc, ghi config từ file json, có các hàm hỗ trợ để lấy giá trị trong cây config.
- `Setting`: là một struct để lưu setting trong database, bao gồm id, key và value.

## Tương lai

Module core trong tương lai sẽ được tách thành một module riêng biệt để có thể dùng trong tất cả các dự án.
