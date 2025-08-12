import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { User } from '../auth/model/User.model';

@Injectable({
  providedIn: 'root'
})
export class StakeholdersService {

  constructor(private http: HttpClient) { }


  getAllUsers() {
  return this.http.get<User[]>('http://localhost:8080/users/all');
}

}
