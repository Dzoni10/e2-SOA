import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class NotificationsService {
  private ws!: WebSocket;
  private _notifications = new Subject<any>();

  public notifications$ = this._notifications.asObservable();

  connect(userId: string) {
    this.ws = new WebSocket(`ws://localhost:8070/followers/ws?userId=${userId}`);

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      this._notifications.next(data);
    };

    this.ws.onclose = () => {
      console.log('WebSocket closed, reconnecting in 5s...');
      setTimeout(() => this.connect(userId), 5000);
    };
  }

  disconnect() {
    this.ws.close();
  }
}
