package com.example.springweb.adapters.userservice;

import com.example.springweb.entities.UserEntity;
import com.example.springweb.services.repositories.UserRepository;
import io.grpc.stub.StreamObserver;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class UserServiceGRPCController extends UserServiceGrpc.UserServiceImplBase {

  public UserRepository userRepository;

  public UserRepository getUserRepository() {
    return userRepository;
  }

  @Autowired
  public void setUserRepository(UserRepository userRepository) {
    this.userRepository = userRepository;
  }

  @Override
  public void getProfile(GetProfileRequest request, StreamObserver<GetProfileResponse> responseObserver) {
   long userId = request.getUserId();
    UserEntity userById = userRepository.findUserById(userId);
    GetProfileResponse getProfileResponse = GetProfileResponse.newBuilder().setId(userById.getId()).setEmail(userById.getEmail()).build();
    responseObserver.onNext(getProfileResponse);
    responseObserver.onCompleted();
  }
}
