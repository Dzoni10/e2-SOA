import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface FollowRequest {
  followerId: number;
  followingId: number;
}

export interface FollowResponse {
  success: boolean;
  message: string;
  following?: boolean;
}

export interface User {
  userId: number;
  username?: string;
  name?: string;
}

export interface UserFollowStats {
  userId: number;
  followersCount: number;
  followingCount: number;
}

export interface FollowRecommendation {
  recommendedUser: User;
  mutualFollowers: string[];
  recommendationReason: string;
}

export interface RecommendationsResponse {
  recommendations: FollowRecommendation[];
  totalCount: number;
}

@Injectable({
  providedIn: 'root'
})
export class FollowerService {
  private apiUrl = 'http://localhost:8070/followers';

  constructor(private http: HttpClient) { }

  followUser(followerId: number, followingId: number): Observable<FollowResponse> {
    const request: FollowRequest = { followerId, followingId };
    return this.http.post<FollowResponse>(`${this.apiUrl}/follow`, request);
  }

  unfollowUser(followerId: number, followingId: number): Observable<FollowResponse> {
    const request: FollowRequest = { followerId, followingId };
    return this.http.post<FollowResponse>(`${this.apiUrl}/unfollow`, request);
  }

  isFollowing(followerId: number, followingId: number): Observable<{isFollowing: boolean}> {
    return this.http.get<{isFollowing: boolean}>(`${this.apiUrl}/is-following/${followerId}/${followingId}`);
  }

  getFollowers(userId: number): Observable<User[]> {
    return this.http.get<User[]>(`${this.apiUrl}/users/${userId}/followers`);
  }

  getFollowing(userId: number): Observable<User[]> {
    return this.http.get<User[]>(`${this.apiUrl}/users/${userId}/following`);
  }

  getUserStats(userId: number): Observable<UserFollowStats> {
    return this.http.get<UserFollowStats>(`${this.apiUrl}/users/${userId}/stats`);
  }

  canUserComment(commenterId: number, authorId: number): Observable<{canComment: boolean}> {
    return this.http.get<{canComment: boolean}>(`${this.apiUrl}/can-comment/${commenterId}/${authorId}`);
  }

  createUser(userId: number, name: string, username?: string): Observable<FollowResponse> {
    return this.http.post<FollowResponse>(`${this.apiUrl}/users`, { userId, name, username });
  }

  getFollowRecommendations(userId: number, limit: number = 10): Observable<RecommendationsResponse> {
    return this.http.get<RecommendationsResponse>(`${this.apiUrl}/users/${userId}/recommendations?limit=${limit}`);
  }
}