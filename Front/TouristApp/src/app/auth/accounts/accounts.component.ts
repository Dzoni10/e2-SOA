import { Component, OnInit } from '@angular/core';
import { Account } from '../model/Account.model';
import { AuthService } from '../auth.service';
import { StakeholdersService } from 'src/app/stakeholders/stakeholders.service';

@Component({
  selector: 'app-accounts',
  templateUrl: './accounts.component.html',
  styleUrls: ['./accounts.component.css']
})
export class AccountsComponent implements OnInit {

  accounts: Account[]=[];

  loading=true;

  constructor(private authService: AuthService, private stakeholderService: StakeholdersService){}

  ngOnInit(): void {
    this.authService.getAccounts().subscribe({
      next: (data) => {
        this.accounts = data;
        this.loading = false;
      },
      error: (err) => {
        console.error('Error during loading accounts:', err);
        this.loading = false;
      }
    });
  }

  
blockUser(id: number){
  this.stakeholderService.blockUser(id).subscribe();
}
  

}
