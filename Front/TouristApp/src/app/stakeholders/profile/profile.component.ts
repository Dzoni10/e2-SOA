import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { UserProfile } from '../model/user.model';
import { StakeholdersService } from '../stakeholders.service';
import { AuthService } from 'src/app/auth/auth.service';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent {
  profileForm: FormGroup;
  isEditing = false;
  loading = true;
  saving = false;
  successMessage = '';
  errorMessage = '';
  userId: number = 0; // You'll need to get this from your auth service or route params

  constructor(
    private fb: FormBuilder,
    private stakeholderService: StakeholdersService,
    private authService: AuthService
  ) {
    this.profileForm = this.fb.group({
      id: [''],
      name: ['', [Validators.required, Validators.maxLength(20)]],
      surname: ['', [Validators.required, Validators.maxLength(40)]],
      username: [{value: '', disabled: true}],
      role: [{value: '', disabled: true}],
      image: [''],
      bio: [''],
      moto: ['']
    });
  }

  ngOnInit(): void {
    this.userId = this.authService.getCurrentUser()?.userId!
    this.loadUserProfile();
  }

  loadUserProfile(): void {
    this.loading = true;
    this.clearMessages();

    this.stakeholderService.getUserProfile(this.userId).subscribe({
      next: (profile: any) => {
        this.profileForm.patchValue(profile);
        this.loading = false;
      },
      error: (error) => {
        this.errorMessage = 'Failed to load profile';
        this.loading = false;
        console.error('Error loading profile:', error);
      }
    });
  }

  startEdit(): void {
    this.isEditing = true;
    this.clearMessages();
  }

  cancelEdit(): void {
    this.isEditing = false;
    this.loadUserProfile(); // Reload original data
    this.clearMessages();
  }

  onSubmit(): void {
    if (this.profileForm.invalid || this.saving) {
      return;
    }

    this.saving = true;
    this.clearMessages();

    // Get the form values and create UserProfile object
    const profileData: UserProfile = {
      id: this.profileForm.get('id')?.value,
      name: this.profileForm.get('name')?.value,
      surname: this.profileForm.get('surname')?.value,
      username: this.profileForm.get('username')?.value,
      role: this.profileForm.get('role')?.value,
      image: this.profileForm.get('image')?.value,
      bio: this.profileForm.get('bio')?.value,
      moto: this.profileForm.get('moto')?.value
    };

    this.stakeholderService.editUserProfile(profileData).subscribe({
      next: (response) => {
        this.successMessage = 'Profile updated successfully!';
        this.isEditing = false;
        this.saving = false;
        this.profileForm.markAsPristine();
        
        // Clear success message after 3 seconds
        setTimeout(() => {
          this.successMessage = '';
        }, 3000);
      },
      error: (error) => {
        this.errorMessage = 'Failed to update profile. Please try again.';
        this.saving = false;
        console.error('Error updating profile:', error);
      }
    });
  }

  getRoleClass(): string {
    const role = this.profileForm.get('role')?.value?.toLowerCase();
    return role || '';
  }

  private clearMessages(): void {
    this.successMessage = '';
    this.errorMessage = '';
  }
}
