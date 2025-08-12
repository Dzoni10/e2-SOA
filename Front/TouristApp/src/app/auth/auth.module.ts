import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {MatFormFieldModule} from '@angular/material/form-field'
import { ReactiveFormsModule } from '@angular/forms';
import {MatInputModule} from '@angular/material/input'
import {MatButtonModule} from '@angular/material/button'
import { LoginComponent } from './login/login.component';
import { SignupComponent } from './signup/signup.component';
import {MatSelectModule} from '@angular/material/select';
import {MatSnackBarModule} from '@angular/material/snack-bar';
import { RouterModule } from '@angular/router';
import { AccountsComponent } from './accounts/accounts.component';
import { MatCardModule } from '@angular/material/card';
import {MatGridListModule} from '@angular/material/grid-list';

@NgModule({
  declarations: [
    LoginComponent,
    SignupComponent,
    AccountsComponent
  ],
  imports: [
    CommonModule,
    MatFormFieldModule,
    ReactiveFormsModule,
    MatInputModule,
    MatButtonModule,
    MatSelectModule,
    RouterModule,
    MatCardModule,
    MatGridListModule
  ],
  exports:[
    
  ]
})
export class AuthModule { }
