import { Component } from '@angular/core';
import { FormGroup,FormControl,Validators } from '@angular/forms';
import { AuthService } from '../auth.service';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';
import { trigger, state, style, transition, animate } from '@angular/animations';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
  animations: [
    trigger('slideIn', [
      state('void', style({ transform: 'translateY(0)', opacity: 0 })),
      transition(':enter', [
        animate('2.5s ease-out', style({ transform: 'translateY(0)', opacity: 1 }))
      ])
    ])
  ]
})
export class LoginComponent {

  loginForm = new FormGroup({
    username: new FormControl('', [Validators.required]),
    password: new FormControl('',[Validators.required])
  });

  constructor(private authService:AuthService, private router: Router, private snackBar:MatSnackBar){}


  login():void{
    if(this.loginForm.value.username !== null){
    this.authService.login(this.loginForm.get('username')!.value!,this.loginForm.get('password')!.value!).subscribe({
      next: (res)=>{
        this.authService.saveToken(res.token);
        console.log('Logged in succesfully');
        this.snackBar.open("Logged successfully!","Close",{duration:3000,horizontalPosition:"center"})
        const user = this.authService.getCurrentUser();
        if(user?.role===0){
          this.router.navigate(['/userAccounts']);
        }
        else if(user?.role===1){

          this.router.navigate(['/toursList']);
        }
      },
      error: (err)=>{
        console.error('Login failed',err);
        this.snackBar.open("Cannot login","Close",{duration:3000,horizontalPosition:"center"})
      }
    });
  }
  }
}
