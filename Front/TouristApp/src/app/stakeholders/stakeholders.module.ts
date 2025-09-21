import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { UserAccountsComponent } from './user-accounts/user-accounts.component';
import { ProfileComponent } from './profile/profile.component';
import { ReactiveFormsModule } from '@angular/forms';

@NgModule({
  declarations: [
    
  
    ProfileComponent
  ],
  imports: [
    CommonModule,
    MatCardModule,
    BrowserAnimationsModule,
    ReactiveFormsModule
  ]
})
export class StakeholdersModule { }
