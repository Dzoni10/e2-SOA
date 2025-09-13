import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { Position } from '../model/position.model';
import { AuthService } from 'src/app/auth/auth.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ToursService } from '../tours.service';
import * as L from 'leaflet';

@Component({
  selector: 'app-position',
  templateUrl: './position.component.html',
  styleUrls: ['./position.component.css']
})
export class PositionComponent implements OnInit {
  position!: Position;
  map!: L.Map;
  marker!: L.Marker;
  //user!: User;
  userId!: number
  hasExistingPosition = false;
  
  @ViewChild('mapContainer') mapContainer!: ElementRef;

  constructor(
    private authService: AuthService, 
    private tourService: ToursService, 
    private snackBar: MatSnackBar
  ) {}

  ngOnInit() {
    const user = this.authService.getCurrentUser();
    
    if (!user || !user.userId) {
      this.snackBar.open("You must be logged in to view current position", "Close", {
        duration: 3000, 
        horizontalPosition: "center"
      });
      return;
    }
    
    this.userId = user.userId;
    this.ngAfterViewInit()
  }

  ngAfterViewInit(){
    this.loadPosition(this.userId);
  }

  async loadPosition(userId: number) {
    this.tourService.getPosition(userId).subscribe({
      next: (position) => {
        this.position = position;
        this.hasExistingPosition = true;
        console.log("djordje")
        this.initMap();
        console.log("ristic")
        this.addMarker(position.latitude, position.longitude);
      },
      error: (error) => {
        console.log('No existing position found, will create new one on first save');
        this.hasExistingPosition = false;
        // Initialize with default position (e.g., Belgrade, Serbia)
        this.position = {
          userId: userId,
          latitude: 44.817,
          longitude: 20.456
        };
        this.initMap();
      }
    });
  }

  initMap() {
    console.log("vil")
    if (!this.mapContainer) return;
    console.log("ker")
    // Initialize map with current position or default location
    this.map = L.map(this.mapContainer.nativeElement).setView(
      [this.position.latitude, this.position.longitude], 
      13
    );

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.map);

    // Add click listener to update position
    this.map.on('click', (e: L.LeafletMouseEvent) => {
      this.updateMapPosition(e.latlng.lat, e.latlng.lng);
    });
  }

  addMarker(lat: number, lng: number) {
    // Remove existing marker if it exists
    if (this.marker) {
      this.map.removeLayer(this.marker);
    }

    // Create new marker
    this.marker = L.marker([lat, lng], { draggable: true }).addTo(this.map);
    
    // Handle marker drag
    this.marker.on('dragend', (event: any) => {
      const latlng = event.target.getLatLng();
      this.updateMapPosition(latlng.lat, latlng.lng);
    });
  }

  updateMapPosition(lat: number, lng: number) {
    this.position.latitude = lat;
    this.position.longitude = lng;
    
    // Add or update marker
    this.addMarker(lat, lng);
    
    // Center map on new position
    this.map.setView([lat, lng]);
  }

  savePosition() {
    if (!this.userId) {
      this.snackBar.open("User not found", "Close", { duration: 3000 });
      return;
    }

    if (this.hasExistingPosition) {
      // Update existing position
      this.tourService.updatePosition(this.userId, this.position).subscribe({
        next: (updatedPosition) => {
          this.position = updatedPosition;
          this.snackBar.open("Position updated successfully", "Close", { 
            duration: 3000,
            horizontalPosition: "center"
          });
        },
        error: (error) => {
          this.snackBar.open("Failed to update position", "Close", { 
            duration: 3000,
            horizontalPosition: "center"
          });
          console.error('Update position error:', error);
        }
      });
    } else {
      // Create new position
      this.tourService.initializePosition(this.userId, this.position).subscribe({
        next: (createdPosition) => {
          this.position = createdPosition;
          this.hasExistingPosition = true;
          this.snackBar.open("Position created successfully", "Close", { 
            duration: 3000,
            horizontalPosition: "center"
          });
        },
        error: (error) => {
          this.snackBar.open("Failed to create position", "Close", { 
            duration: 3000,
            horizontalPosition: "center"
          });
          console.error('Create position error:', error);
        }
      });
    }
  }
}