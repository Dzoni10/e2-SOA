import { Component, OnInit } from '@angular/core';
import { User } from 'src/app/auth/model/User.model';
import { StakeholdersService } from '../stakeholders.service';
import { MatCard } from '@angular/material/card';


@Component({
  selector: 'app-user-accounts',
  templateUrl: './user-accounts.component.html',
  styleUrls: ['./user-accounts.component.css']
})
export class UserAccountsComponent implements OnInit {

  users: User[]=[];


  constructor(private stakeholderService: StakeholdersService){}

  ngOnInit(): void {
    this.stakeholderService.getAllUsers().subscribe({
      next:(data) => this.users=data,
      error:(err)=> console.error('Error loading users',err)
    })
  }

getRoleName(role: number): string{
  switch (role) {
      case 0: return 'Admin';
      case 1: return 'Guide';
      case 2: return 'Tourist';
      default: return 'Unknown';
    }
}


}
