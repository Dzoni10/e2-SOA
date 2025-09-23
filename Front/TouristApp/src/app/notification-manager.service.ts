import { Injectable } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';

import { Subscription } from 'rxjs';
import { NotificationsService } from './stakeholders/notifications.service';
import { AuthService } from './auth/auth.service';

@Injectable({
  providedIn: 'root'
})
export class NotificationManagerService {
  private notificationSubscription?: Subscription;
  private isInitialized = false;

  constructor(
    private notificationService: NotificationsService,
    private snackBar: MatSnackBar,
    private authService: AuthService
  ) {}

  /**
   * Inicijalizuje globalne notifikacije - poziva se iz AppComponent
   */
  initializeGlobalNotifications(): void {
    if (this.isInitialized) {
      return; // Već inicijalizovano
    }

    const currentUser = this.authService.getCurrentUser();
    if (!currentUser?.userId) {
      console.warn('User not authenticated, cannot initialize notifications');
      return;
    }

    // Konektuj se na notification service
    this.notificationService.connect(currentUser.userId.toString());

    // Pretplati se na notifikacije
    this.notificationSubscription = this.notificationService.notifications$.subscribe({
      next: (notification) => {
        console.log("jesam")
        this.showNotificationSnackbar(notification);
        console.log("jesam2")
      },
      error: (error) => {
        console.error('Error receiving notification:', error);
      }
    });

    this.isInitialized = true;
    console.log('Global notifications initialized for user:', currentUser.userId);
  }

  /**
   * Prikazuje notifikaciju u snackbar-u
   */
  private showNotificationSnackbar(notification: any): void {
    console.log("u notif sam")
    const message = `New blog by author ${notification.authorId}: ${notification.title}`;
    
    this.snackBar.open("message", 'Close', {
      duration: 9000,
      horizontalPosition: 'center',
      //verticalPosition: 'bottom',
      panelClass: ['notification-snackbar'] // Optional: za custom styling
    });

    console.log('Notification displayed:', notification);
  }

  /**
   * Čisti resurse - poziva se kada se user logout-uje ili aplikacija se gasi
   */
  cleanup(): void {
    if (this.notificationSubscription) {
      this.notificationSubscription.unsubscribe();
      this.notificationSubscription = undefined;
    }

    this.notificationService.forceDisconnect(); // Koristi forceDisconnect umesto disconnect
    this.isInitialized = false;
    console.log('Global notifications cleaned up');
  }

  /**
   * Restartuje notifikacije - korisno nakon login/logout
   */
  restart(): void {
    this.cleanup();
    this.initializeGlobalNotifications();
  }
}