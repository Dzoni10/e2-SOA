import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { window } from 'rxjs';
import { AuthService } from 'src/app/auth/auth.service';
import { DecodedToken } from 'src/app/auth/model/decodedToken';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css']
})
export class NavbarComponent implements OnInit {

  user: DecodedToken | null=null;
  private subscription!: Subscription;

  constructor(private authService:AuthService, private router: Router){}

  ngOnInit(): void {
    this.subscription = this.authService.currentUser$.subscribe(user=>{
      this.user=user;
    });
  }

  isLoggedIn(): boolean{
    return this.user!==null;
  }

  logout(): void{
    this.authService.logout();
    this.router.navigate(['/login']);
  }

  ngOnDestroy(): void {
    this.subscription.unsubscribe();
  }

  isLoggedAdmin(): boolean{
    return this.user?.role===0
  }

  isLoggedAuthor(): boolean{
    return this.user?.role===1
  }
}
