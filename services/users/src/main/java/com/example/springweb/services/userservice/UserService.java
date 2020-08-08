package com.example.springweb.services.userservice;

import com.example.springweb.entities.UserEntity;
import com.example.springweb.services.repositories.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import java.security.SecureRandom;

@Service
public class UserService {

  private UserRepository userRepository;
  private SHA256PasswordEncoder passwordEncoder;
  private JwtEncoder jwtEncoder;

  @Autowired
  public void setUserRepository(UserRepository userRepository) {
    this.userRepository = userRepository;
  }

  @Autowired
  public void setPasswordEncoder(SHA256PasswordEncoder passwordEncoder) {
    this.passwordEncoder = passwordEncoder;
  }

  @Autowired
  public void setJwtEncoder(JwtEncoder jwtEncoder) {
    this.jwtEncoder = jwtEncoder;
  }

  public UserEntity registerUser(String email, String password) {
    UserEntity userEntity = new UserEntity();
    SecureRandom random = new SecureRandom();
    byte[] baseSalt = new byte[16];
    random.nextBytes(baseSalt);

    byte[] salt = random.toString().getBytes();
    byte[] encodedPassword = passwordEncoder.encode(password, salt);
    userEntity.setEmail(email);
    userEntity.setHash(encodedPassword);
    userEntity.setSalt(salt);
    return userRepository.saveUser(userEntity);
  }

  public JwtTokenDto auth(String email, String password) {
    UserEntity user = userRepository.findUserByEmail(email);
    if (user == null) {
      throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "Wrong email or password");
    }

    if (!passwordEncoder.match(user.getHash(), password, user.getSalt())) {
      throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "Wrong email or password");
    }
    return jwtEncoder.generateToken(user);
  }

  public JwtTokenDto refreshToken(String token) {

    JwtClaim decodeToken = jwtEncoder.decodeToken(token);
    if (decodeToken == null){
      throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "Could not decode token, or it is expired");
    }

    UserEntity userByEmail = userRepository.findUserByEmail(decodeToken.email);
    if (userByEmail == null) {
      throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "No user found for this jwt token");
    }
    return jwtEncoder.generateToken(userByEmail);
  }
}
