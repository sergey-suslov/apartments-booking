package com.example.springweb.adapters;

import com.example.springweb.adapters.dtos.RegisterUserDto;
import com.example.springweb.entities.UserEntity;
import com.example.springweb.services.userservice.JwtTokenDto;
import com.example.springweb.services.userservice.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

import javax.validation.Valid;

@RestController
public class UserController {
  UserService userService;

  @Autowired
  public void setUserService(UserService userService) {
    this.userService = userService;
  }

  @PostMapping("/")
  @ResponseBody
  UserEntity register(@Valid @RequestBody RegisterUserDto registerUserDto) {
    return userService.registerUser(registerUserDto.getEmail(), registerUserDto.getPassword());
  }

  @PostMapping("/auth")
  @ResponseBody
  JwtTokenDto auth(@Valid @RequestBody RegisterUserDto registerUserDto) {
    return userService.auth(registerUserDto.getEmail(), registerUserDto.getPassword());
  }

  @PostMapping("/refresh")
  @ResponseBody
  JwtTokenDto refresh(@RequestBody String token) {
    return userService.refreshToken(token);
  }
}
