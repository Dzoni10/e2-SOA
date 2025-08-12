import { Injectable } from '@angular/core';
import {HttpClient, HttpStatusCode} from '@angular/common/http'
import { Router } from '@angular/router';
import { Login } from './model/login.model';
import { Register } from './model/register.model';
import { Observable } from 'rxjs';
import { enviroment } from 'src/env/enviroment';
import { LocalizedString } from '@angular/compiler';
import { jwtDecode } from 'jwt-decode';
import { DecodedToken } from './model/decodedToken';
import { BehaviorSubject } from 'rxjs';
import { Account } from './model/Account.model';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  private currentUserSubject = new BehaviorSubject<DecodedToken | null>(this.loadUserFromToken());
  currentUser$ = this.currentUserSubject.asObservable();

  private apiUrl = 'http://localhost:8080/users'
  constructor(private http: HttpClient) { }

  register(registration: Register): Observable<any>{
    return this.http.post('http://localhost:8080/users',registration);
  }

  login(username: string, password: string){
    const body = {username, password};
    return this.http.post<{token:string}>(`${this.apiUrl}/login`,body)
  }

  getAccounts(): Observable<Account[]>{
    return this.http.get<Account[]>(`${this.apiUrl}/all`)
  }

  saveToken(token: string){
    localStorage.setItem('jwtToken',token)
    const decoded=jwtDecode<DecodedToken>(token);
    this.currentUserSubject.next(decoded);
  }

  getToken():string|null{
    return localStorage.getItem('jwtToken')
  }

  getCurrentUser(): DecodedToken|null{
    return this.currentUserSubject.value;
  }

  logout(){
    localStorage.removeItem('jwtToken');
    this.currentUserSubject.next(null);
  }
  
  private loadUserFromToken(): DecodedToken | null {
    const token = this.getToken();
    if(token){
      try{
        return jwtDecode<DecodedToken>(token);
      }catch{
        return null;
      }
    }
    return null;
  }

}


