import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { User } from '../auth/model/User.model';
import { UserProfile } from './model/user.model';

@Injectable({
  providedIn: 'root'
})
export class StakeholdersService {

  constructor(private http: HttpClient) { }


  getAllUsers() {
  return this.http.get<User[]>('http://localhost:8070/users/all');
}
getUserById(id: number) {
    return this.http.get<User>(`http://localhost:8070/users/${id}`);
  }

getUserProfile(id: number) {
    return this.http.get<UserProfile>(`http://localhost:8070/users/${id}`);
  }

editUserProfile(profile: UserProfile) {
  return this.http.put<UserProfile>(`http://localhost:8070/users/editProfile`, profile)
}

blockUser(id: number) {
  return this.http.get<UserProfile>(`http://localhost:8070/users/${id}/block`);
}
}
