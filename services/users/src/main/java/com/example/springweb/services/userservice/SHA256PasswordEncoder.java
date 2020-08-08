package com.example.springweb.services.userservice;

import lombok.SneakyThrows;
import org.springframework.stereotype.Component;

import javax.crypto.SecretKeyFactory;
import javax.crypto.spec.PBEKeySpec;
import java.security.spec.KeySpec;
import java.util.Arrays;

@Component
public class SHA256PasswordEncoder {
  @SneakyThrows
  public byte[] encode(String rawPassword, byte[] salt) {
    KeySpec spec = new PBEKeySpec(rawPassword.toCharArray(), salt, 100, 128);
    SecretKeyFactory factory = SecretKeyFactory.getInstance("PBKDF2WithHmacSHA1");
    return factory.generateSecret(spec).getEncoded();
  }

  public boolean match(byte[] hash, String rawPassword, byte[] salt) {
    return Arrays.equals(hash, encode(rawPassword, salt));
  }
}
