import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup,Validators } from '@angular/forms';
import { AuthService } from '../auth.service';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: 'app-signup',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.css']
})
export class SignupComponent implements OnInit {

  signupForm!: FormGroup
  constructor(private fb: FormBuilder, private authService: AuthService,private snackBar:MatSnackBar){}


  ngOnInit(): void {
    this.signupForm = this.fb.group({
      name: ['', Validators.required],
      surname: ['', Validators.required],
      username: ['', Validators.required],
      password: ['', Validators.required],
      role: ['', Validators.required]
    });
  }

  register(): void{
    if (this.signupForm.valid){
      this.authService.register(this.signupForm.value).subscribe({
        next: res=> {
          console.log("REGISTER SUCCESS");
          this.snackBar.open("Register successfully!","Close",{duration:3000,horizontalPosition:"center"})
          this.signupForm.reset();
                },
        error: err=>{ 
          console.log("CANNTO REGISTER ",err);
        this.snackBar.open("Cannot register!","Close",{duration:3000,horizontalPosition:"center"})     
    }
    });
  }
}
}
