// 解密函数
import CryptoJS from 'crypto-js';
export function decryptAPIKey(encryptedKeyBase64) {
    // 密钥（与后端一致）
    const key = CryptoJS.enc.Utf8.parse('your-32-byte-encryption-key-here');

    // Base64解码
    const encrypted = CryptoJS.enc.Base64.parse(encryptedKeyBase64);

    // 使用ECB模式解密（最简单）
    const decrypted = CryptoJS.AES.decrypt(
        { ciphertext: encrypted },
        key,
        {
            mode: CryptoJS.mode.ECB, // ECB模式最简单
            padding: CryptoJS.pad.Pkcs7
        }
    );

    // 转换为字符串
    return decrypted.toString(CryptoJS.enc.Utf8);
}
