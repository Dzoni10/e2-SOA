import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import * as L from 'leaflet';
import 'leaflet-routing-machine'; // for road routing
import { Tour, TourStatus, TourStatusInfo } from '../model/tour.model';
import { Keypoint } from '../model/keypoint.model';
import { ToursService } from '../tours.service';
import { forkJoin, Observable } from 'rxjs';
import { MatSnackBar } from '@angular/material/snack-bar';
import { AuthService } from 'src/app/auth/auth.service';
import { FormsModule } from '@angular/forms';
@Component({
  selector: 'app-tour-details',
  templateUrl: './tour-details.component.html',
  styleUrls: ['./tour-details.component.css']
})
export class TourDetailsComponent implements OnInit {
  tour!: Tour;
  keypoints: Keypoint[] = [];
  map!: L.Map;
  @ViewChild('mapContainer') mapContainer!: ElementRef;
  routeControl!: any;
  keypointLayer!: L.LayerGroup
statusInfo!: TourStatusInfo;
canEditStatus = false;
currentUser: any;
TourStatus = TourStatus;
  constructor(
    private route: ActivatedRoute, private tourService: ToursService, private snackBar: MatSnackBar,private authService: AuthService
  ) {}

  ngOnInit(): void {
  const user = this.authService.getCurrentUser();
  this.currentUser = user;

  if (!user || !user.userId) {
    this.snackBar.open("You must be logged in to see tours", "Close", {
      duration: 3000, 
      horizontalPosition: "center"
    });
    return;
  }


    const tourId = this.route.snapshot.paramMap.get('id');
    if (tourId) this.loadTour(tourId);
  }

async loadTour(tourId: string) {
  try {
    const tourResponse = await this.tourService.getTourByID(tourId).toPromise();
    if (tourResponse) {
      this.tour = tourResponse;
    }
    
    const keypointsResponse = await this.tourService.getKeyPointsForTour(tourId).toPromise();
    if (keypointsResponse) {
      this.keypoints = keypointsResponse;
    }

    const statusResponse = await this.tourService.getTourStatusInfo(tourId, this.currentUser.userId).toPromise();
    if (statusResponse) {
      this.statusInfo = statusResponse;
      this.canEditStatus = statusResponse.canEdit;
    }

    this.initMap();
    
  } catch (error) {
    this.showMessage('Failed to load tour details');
  }
}

// ----------------------- MAP -----------------------
  initMap() {
    if (!this.mapContainer) return;

    this.map = L.map(this.mapContainer.nativeElement).setView([44.817, 20.456], 12);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.map);

    this.keypointLayer = L.layerGroup().addTo(this.map);

    this.drawAllKeypoints();
    this.drawRoute();
    this.enableAddingKeypoints();
  }

  private getMarkerIcon(): L.Icon {
    return L.icon({
      iconUrl: 'assets/markinjo.png',
      iconSize: [40, 65],
      iconAnchor: [20, 65],
      popupAnchor: [0, -60]
    });
  }

 // ------------------- CRTA SVE KEYPOINTE -------------------
drawAllKeypoints() {
  if (!this.keypointLayer) this.keypointLayer = L.layerGroup().addTo(this.map);
  this.keypointLayer.clearLayers();

  this.keypoints.forEach(kp => {
    // Marker
    const marker = L.marker([kp.latitude, kp.longitude], {
      draggable: true,
      icon: this.getMarkerIcon(),
      riseOnHover: true
    });

    // Circle za veći clickable area
    const clickCircle = L.circle([kp.latitude, kp.longitude], {
      radius: 30, // u metrima
      color: 'transparent',
      fillColor: 'transparent',
      weight: 0,
      interactive: true
    });

    // Klik na circle otvara popup markera
    clickCircle.on('click', () => marker.openPopup());

    // Popup content
    const popupContent = `
      <div style="min-width:200px">
        <h4>Edit Keypoint</h4>
        <label>Name:</label><br>
        <input id="kp-name-${kp.id}" type="text" value="${kp.name}" style="width:100%"/><br>
        <label>Description:</label><br>
        <textarea id="kp-desc-${kp.id}" style="width:100%">${kp.description}</textarea><br>
        <label>Image:</label><br>
        <input id="kp-img-${kp.id}" type="file" accept="image/*" /><br><br>
        <button id="kp-save-${kp.id}">💾 Save</button>
        <button id="kp-delete-${kp.id}">❌ Delete</button>
      </div>
    `;
    marker.bindPopup(popupContent);

    // Popup open event za listener-e
    marker.on('popupopen', () => {
      const saveBtn = document.getElementById(`kp-save-${kp.id}`);
      if (saveBtn) saveBtn.addEventListener('click', () => this.saveKeypoint(kp, marker), { once: true });

      const delBtn = document.getElementById(`kp-delete-${kp.id}`);
      if (delBtn) delBtn.addEventListener('click', () => this.deleteKeypoint(kp), { once: true });
    });

    // Dragging markera
    marker.on('dragend', (event: any) => {
      const latlng = event.target.getLatLng();
      kp.latitude = latlng.lat;
      kp.longitude = latlng.lng;
      clickCircle.setLatLng(latlng);
      this.drawRoute();
    });

    kp.marker = marker;
    kp.clickCircle = clickCircle;

    this.keypointLayer.addLayer(marker);
    this.keypointLayer.addLayer(clickCircle);
  });
}


  saveKeypoint(kp: Keypoint, marker: L.Marker) {
    const nameInput = (document.getElementById(`kp-name-${kp.id}`) as HTMLInputElement).value;
    const descInput = (document.getElementById(`kp-desc-${kp.id}`) as HTMLTextAreaElement).value;
    const imgInput = (document.getElementById(`kp-img-${kp.id}`) as HTMLInputElement).files?.[0];

    kp.name = nameInput;
    kp.description = descInput;

    const formData = new FormData();
    formData.append("name", kp.name);
    formData.append("description", kp.description || "");
    formData.append("latitude", kp.latitude.toString());
    formData.append("longitude", kp.longitude.toString());
    if (imgInput) formData.append("image", imgInput);

    if (kp.id) {
      this.tourService.updateKeypoint(kp.id, formData).subscribe(() => marker.closePopup());
    }
  }

  deleteKeypoint(kp: Keypoint) {
    if (kp.id) {
        if (kp['marker']) this.keypointLayer.removeLayer(kp['marker']);
        this.keypoints = this.keypoints.filter(k => k !== kp);
        this.drawRoute();

    } else {
      // new unsaved keypoint
      if (kp['marker']) this.keypointLayer.removeLayer(kp['marker']);
      this.keypoints = this.keypoints.filter(k => k !== kp);
      this.drawRoute();
    }
  }

  enableAddingKeypoints() {
  this.map.on('click', (e: L.LeafletMouseEvent) => {
    // Proveri da li je klik unutar bilo kojeg postojeceg circle-a
    const isInsideExisting = this.keypoints.some(kp => {
      if (!kp.clickCircle) return false;
      return kp.clickCircle.getBounds().contains(e.latlng);
    });
    if (isInsideExisting) return; // ne kreiraj novi

    // Novi keypoint
    const newKp: Keypoint = {
      name: 'New Keypoint',
      description: '',
      latitude: e.latlng.lat,
      longitude: e.latlng.lng,
      images: [],
      order: this.keypoints.length
    };

    // Marker
    const marker = L.marker([newKp.latitude, newKp.longitude], {
      draggable: true,
      icon: this.getMarkerIcon(),
      riseOnHover: true
    });

    // Circle
    const clickCircle = L.circle([newKp.latitude, newKp.longitude], {
      radius: 30,
      color: 'transparent',
      fillColor: 'transparent',
      weight: 0,
      interactive: true
    });
    clickCircle.on('click', () => marker.openPopup());

    // Popup content
    const popupContent = `
      <div style="min-width:200px">
        <h4>New Keypoint</h4>
        <label>Name:</label><br>
        <input id="new-kp-name" type="text" style="width:100%"/><br>
        <label>Description:</label><br>
        <textarea id="new-kp-desc" style="width:100%"></textarea><br>
        <label>Image:</label><br>
        <input id="new-kp-img" type="file" accept="image/*" /><br><br>
        <button id="new-kp-save">➕ Add</button>
      </div>
    `;
    marker.bindPopup(popupContent).openPopup();

    marker.on('popupopen', () => {
      const saveBtn = document.getElementById("new-kp-save");
      if (saveBtn) {
        saveBtn.addEventListener("click", () => {
          const nameInput = (document.getElementById("new-kp-name") as HTMLInputElement).value;
          const descInput = (document.getElementById("new-kp-desc") as HTMLTextAreaElement).value;
          const imgInput = (document.getElementById("new-kp-img") as HTMLInputElement).files?.[0];

          newKp.name = nameInput;
          newKp.description = descInput;

          const formData = new FormData();
          formData.append("name", newKp.name);
          formData.append("description", newKp.description || "");
          formData.append("latitude", newKp.latitude.toString());
          formData.append("longitude", newKp.longitude.toString());
          if (imgInput) formData.append("image", imgInput);

          //this.tourService.createKeypoint(formData).subscribe(savedKp => {
            //newKp.id = savedKp.id;
            const tourId = this.route.snapshot.paramMap.get('id');
            newKp.tourId = tourId! 
            marker.closePopup();
            this.drawRoute();
          //});
        }, { once: true });
      }
    });

    // Drag marker
    marker.on('dragend', (event: any) => {
      const latlng = event.target.getLatLng();
      newKp.latitude = latlng.lat;
      newKp.longitude = latlng.lng;
      clickCircle.setLatLng(latlng);
      this.drawRoute();
    });

    newKp.marker = marker;
    newKp.clickCircle = clickCircle;

    this.keypoints.push(newKp);
    this.keypointLayer.addLayer(marker);
    this.keypointLayer.addLayer(clickCircle);
  });
}


  updateKeypointPosition(kp: Keypoint, event: any) {
    const latlng = event.target.getLatLng();
    kp.latitude = latlng.lat;
    kp.longitude = latlng.lng;
    this.keypoints.sort((a, b) => (a.order || 0) - (b.order || 0));
    this.drawRoute();
  }

  drawRoute() {
  if (this.routeControl) this.routeControl.remove();

  const waypoints = this.keypoints
    .filter(kp => kp.latitude && kp.longitude)
    .map(kp => L.Routing.waypoint(L.latLng(kp.latitude, kp.longitude)));

  if (waypoints.length < 2) return;

  this.routeControl = L.Routing.control({
    waypoints,
    router: L.Routing.osrmv1({
      serviceUrl: 'https://router.project-osrm.org/route/v1',
      profile: 'driving'
    }),
    addWaypoints: false,
    //draggableWaypoints: false,
    fitSelectedRoutes: true,
    show: false, // ne prikazuje UI kontrole
  }).addTo(this.map);

  // opcionalno promeni boju linije
  this.routeControl.on('routesfound', (e: any) => {
    const routes = e.routes;
    if (!routes || routes.length === 0) return;

    routes.forEach((r: any) => {
      r.coordinates.forEach((c: any, i: number) => {
        // Leaflet polylines automatski crtaju liniju, boja default je plava
      });
    });
  });
}




publishTour(): void {
  if (!this.canEditStatus) {
    this.showMessage('You are not authorized to change this tour status');
    return;
  }

  this.tourService.publishTour(this.tour.id!,this.currentUser.userId).subscribe({
    next: () => {
      this.tour.status = TourStatus.Published;
      this.showMessage('Tour published successfully');
      this.refreshStatusInfo();
    },
    error: (error) => {
      console.error('Failed to publish tour:', error);
      this.showMessage('Failed to publish tour: ' + (error.error?.error || 'Unknown error'));
    }
  });
}

archiveTour(): void {
  if (!this.canEditStatus) {
    this.showMessage('You are not authorized to change this tour status');
    return;
  }

  this.tourService.archiveTour(this.tour.id!, this.currentUser.userId).subscribe({
    next: () => {
      this.tour.status = TourStatus.Archived;
      this.showMessage('Tour archived successfully');
      this.refreshStatusInfo();
    },
    error: (error) => {
      console.error('Failed to archive tour:', error);
      this.showMessage('Failed to archive tour: ' + (error.error?.error || 'Unknown error'));
    }
  });
}

reactivateTour(): void {
  if (!this.canEditStatus) {
    this.showMessage('You are not authorized to change this tour status');
    return;
  }

  this.tourService.reactivateTour(this.tour.id!, this.currentUser.userId).subscribe({
    next: () => {
      this.tour.status = TourStatus.Published;
      this.showMessage('Tour reactivated successfully');
      this.refreshStatusInfo();
    },
    error: (error) => {
      console.error('Failed to reactivate tour:', error);
      this.showMessage('Failed to reactivate tour: ' + (error.error?.error || 'Unknown error'));
    }
  });
}

private refreshStatusInfo(): void {
  if (this.tour.id) {
    this.tourService.getTourStatusInfo(this.tour.id, this.currentUser.userId).subscribe({
      next: (statusInfo) => {
        this.statusInfo = statusInfo;
      },
      error: (error) => {
        console.error('Failed to refresh status info:', error);
      }
    });
  }
}

private showMessage(message: string): void {
  this.snackBar.open(message, 'Close', {
    duration: 3000,
    horizontalPosition: 'center',
    verticalPosition: 'bottom'
  });
}

getStatusDisplayName(): string {
  if (this.statusInfo?.currentStatus) return this.statusInfo.currentStatus;
  return 'Unknown';
}

canPublish(): boolean {
  return this.canEditStatus && (this.statusInfo?.currentStatus === 'Draft');
}

canArchive(): boolean {
  return this.canEditStatus && (this.statusInfo?.currentStatus === 'Published');
}

canReactivate(): boolean {
  return this.canEditStatus && (this.statusInfo?.currentStatus === 'Archived');
}

  drawKeypoints() {
    // Create layer group if not exists
  if (!this.keypointLayer) {
    this.keypointLayer = L.layerGroup().addTo(this.map);
  } else {
    this.keypointLayer.clearLayers(); // remove old markers
  }

  this.keypoints.forEach(kp => {
    const bigIcon = L.icon({
      iconUrl: 'assets/markinjo.png',   // use default or custom
      iconSize: [40, 65],                  // bigger marker
      iconAnchor: [20, 65],                // so it "sits" on point
      popupAnchor: [0, -60],
    });
    const marker = L.marker([kp.latitude, kp.longitude], {
      draggable: true,
      icon: bigIcon,
      riseOnHover: true // makes marker always stay "on top"
    });//.addTo(this.map);
    this.keypointLayer.addLayer(marker); // ✅ only add to layer group
    // HTML form inside popup
    const popupContent = `
      <div style="min-width:200px">
        <h4>Edit Keypoint</h4>
        <label>Name:</label><br>
        <input id="kp-name-${kp.id}" type="text" value="${kp.name}" style="width:100%"/><br>
        
        <label>Description:</label><br>
        <textarea id="kp-desc-${kp.id}" style="width:100%">${kp.description}</textarea><br>
        
        <label>Image:</label><br>
        <input id="kp-img-${kp.id}" type="file" accept="image/*" /><br><br>
        
        <button id="kp-save-${kp.id}" style="margin-top:5px">💾 Save</button>
      </div>
    `;

    marker.bindPopup(popupContent);

    marker.on('popupopen', () => {
      const saveBtn = document.getElementById(`kp-save-${kp.id}`);
      if (saveBtn) {
        saveBtn.addEventListener('click', () => {
          const nameInput = (document.getElementById(`kp-name-${kp.id}`) as HTMLInputElement).value;
          const descInput = (document.getElementById(`kp-desc-${kp.id}`) as HTMLTextAreaElement).value;
          const imgInput = (document.getElementById(`kp-img-${kp.id}`) as HTMLInputElement).files?.[0];

          kp.name = nameInput;
          kp.description = descInput;

          // build form data for backend (if image upload is needed)
          const formData = new FormData();
          formData.append("name", kp.name);
          formData.append("description", kp.description || "");
          formData.append("latitude", kp.latitude.toString());
          formData.append("longitude", kp.longitude.toString());
          if (imgInput) formData.append("image", imgInput);

          this.tourService.updateKeypoint(kp.id!, formData).subscribe(() => {
            marker.closePopup();
          });
        });
      }
    });

    marker.on('dragend', (event: any) => this.updateKeypointPosition(kp, event));
    kp['marker'] = marker;
  });

  // ✅ clear separation: add new keypoint only when clicking on the map background
  this.map.on('click', (event: L.LeafletMouseEvent) => {
    const latlng = event.latlng;
    this.createKeypoint(latlng.lat, latlng.lng);
  });
}

createKeypoint(lat: number, lng: number) {
  const newKeypoint: Keypoint = {
    // id ne definišemo – backend će ga dodeliti tek kada se snimi
    name: "",
    description: "",
    latitude: lat,
    longitude: lng,
    images: [],
    tourId: this.route.snapshot.paramMap.get("id")! // odmah vežemo za turu
  };

  const marker = L.marker([lat, lng], { draggable: true }).addTo(this.map);

  const popupContent = `
    <div style="min-width:200px">
      <h4>New Keypoint</h4>
      <label>Name:</label><br>
      <input id="new-kp-name" type="text" style="width:100%"/><br>
      
      <label>Description:</label><br>
      <textarea id="new-kp-desc" style="width:100%"></textarea><br>
      
      <label>Image:</label><br>
      <input id="new-kp-img" type="file" accept="image/*" /><br><br>
      
      <button id="new-kp-save">➕ Add</button>
    </div>
  `;

  marker.bindPopup(popupContent).openPopup();

  marker.on("popupopen", () => {
    const saveBtn = document.getElementById("new-kp-save");
    if (saveBtn) {
      saveBtn.addEventListener("click", () => {
        const nameInput = (document.getElementById("new-kp-name") as HTMLInputElement).value;
        const descInput = (document.getElementById("new-kp-desc") as HTMLTextAreaElement).value;
        const imgInput = (document.getElementById("new-kp-img") as HTMLInputElement).files?.[0];

        newKeypoint.name = nameInput;
        newKeypoint.description = descInput;

        if (imgInput) {
          //newKeypoint.images = [imgInput]; // lokalno čuvamo fajl
        }

        // Samo dodajemo u front niz, bez poziva backend-a
        this.keypoints.push(newKeypoint);

        marker.closePopup();
      });
    }
  });

  marker.on("dragend", (event: any) => this.updateKeypointPosition(newKeypoint, event));

  // Čuvamo marker u objektu radi kasnijeg manipulisanja
  (newKeypoint as any).marker = marker;
}


  saveTourUpdates() {
  if (!this.tour) return;

  const updateCalls: Observable<any>[] = [];

  this.keypoints.forEach((kp, index) => {
    const formData = new FormData();
    formData.append("name", kp.name);
    formData.append("description", kp.description || "");
    formData.append("latitude", kp.latitude.toString());
    formData.append("longitude", kp.longitude.toString());
    formData.append("order", (index + 1).toString()); // redosled iz niza
    formData.append("tourId", this.tour!.id!);

    if (kp.images && kp.images.length > 0) {
      // Ako si čuvao fajlove u kp.images
      kp.images.forEach(img => formData.append("images", img));
    }

    if (kp.id) {
      // Update postojeće tačke
      updateCalls.push(this.tourService.updateKeypoint(kp.id, formData));
    } else {
      // Novi keypoint → create
      updateCalls.push(this.tourService.createKeypoint(formData));
    }
  });

  forkJoin(updateCalls).subscribe({
    next: (results) => {
      console.log("All keypoints saved in order:", results);

      // Onda update dužine ture
      this.tourService.updateTourLength(this.tour!.id!).subscribe({
        next: () => console.log("Tour length updated"),
        error: (err) => console.error("Tour length update failed", err)
      });
    },
    error: (err) => console.error("Keypoint save failed", err)
  });
}


updateTourCost(): void {
  if (!this.tour || !this.canEditStatus) return;

  this.tourService.updateTourCost(this.tour.id!, this.tour.cost).subscribe({
    next: () => {
      this.showMessage('Tour cost updated successfully');
    },
    error: (error) => {
      console.error('Failed to update cost:', error);
      this.showMessage('Failed to update cost');
    }
  });
}

}

