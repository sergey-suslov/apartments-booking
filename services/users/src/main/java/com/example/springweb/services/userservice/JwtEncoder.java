package com.example.springweb.services.userservice;

import com.example.springweb.entities.UserEntity;
import io.jsonwebtoken.Jwt;
import io.jsonwebtoken.SignatureAlgorithm;
import io.jsonwebtoken.impl.DefaultClaims;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import io.jsonwebtoken.Jwts;

import java.util.Date;
import java.util.HashMap;

@Component
public class JwtEncoder {

  final private int JWT_TOKEN_VALIDITY = 5 * 60;
  @Value("${jwt.secret}")
  String jwtSecret;

  public JwtTokenDto generateToken(UserEntity userEntity)  {
    HashMap<String, Object> claim = new HashMap<>();
    claim.put("id", userEntity.getId().toString());
    claim.put("email", userEntity.getEmail());

    Date expiration = new Date(System.currentTimeMillis() + JWT_TOKEN_VALIDITY * 1000);
    JwtTokenDto jwtTokenDto = new JwtTokenDto();
    jwtTokenDto.setExpirationDate(expiration);
    jwtTokenDto.setJwt(Jwts.builder().setClaims(claim).setIssuedAt(new Date(System.currentTimeMillis()))
        .setExpiration(expiration)
        .signWith(SignatureAlgorithm.HS512, this.jwtSecret).compact());
    return jwtTokenDto;
  }

  public JwtClaim decodeToken(String token) {
    try {
      DefaultClaims parse = (DefaultClaims) Jwts.parser().setSigningKey(jwtSecret).parse(token).getBody();
      JwtClaim jwtClaim = new JwtClaim();
      jwtClaim.id = Long.valueOf(parse.get("id", String.class));
      jwtClaim.email = parse.get("email", String.class);
      return jwtClaim;
    } catch (Exception e) {
      return null;
    }
  }
}
